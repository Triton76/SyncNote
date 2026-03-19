package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"SyncNote/syncnote/rpc/internal/model"
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GrantPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGrantPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantPermissionLogic {
	return &GrantPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Permission Management (对应 note_permissions 表) ---
func (l *GrantPermissionLogic) GrantPermission(in *syncnoterpc.GrantPermissionReq) (*syncnoterpc.PermissionResp, error) {
	if in == nil {
		return nil, errors.New("request is nil")
	}
	if in.NoteId == "" || in.OperatorId == "" {
		return nil, errors.New("noteId and operatorId are required")
	}
	if in.TargetTeamId != "" && (in.TargetUserId != "" || in.TargetUserEmail != "") {
		return nil, errors.New("targetTeamId and targetUser* must be exactly one")
	}
	if in.TargetTeamId == "" && in.TargetUserId == "" && in.TargetUserEmail == "" {
		return nil, errors.New("targetUserId/targetUserEmail/targetTeamId is required")
	}
	if in.TargetUserId != "" && in.TargetUserEmail != "" {
		return nil, errors.New("targetUserId and targetUserEmail cannot both be set")
	}

	resolvedUserID := strings.TrimSpace(in.TargetUserId)
	if resolvedUserID == "" && strings.TrimSpace(in.TargetUserEmail) != "" {
		var resolveErr error
		resolvedUserID, resolveErr = l.resolveUserIDByEmail(strings.TrimSpace(in.TargetUserEmail))
		if resolveErr != nil {
			return nil, resolveErr
		}
	}

	role, ok := roleToDB(in.Role)
	if !ok {
		return nil, errors.New("invalid role")
	}
	if role == "owner" {
		return nil, errors.New("cannot grant owner role")
	}

	operator, err := l.svcCtx.NotePermissionsModel.FindOneByNoteIdUserId(l.ctx, in.NoteId, sql.NullString{String: in.OperatorId, Valid: true})
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("permission denied")
		}
		return nil, err
	}
	operatorRole := strings.ToLower(operator.Role)
	if operatorRole != "owner" && operatorRole != "admin" {
		return nil, errors.New("permission denied")
	}

	now := time.Now().UnixMilli()
	status := "active"
	permission := &model.NotePermissions{
		NoteId:    in.NoteId,
		GrantedBy: in.OperatorId,
		Role:      role,
		Status:    status,
		GrantedAt: now,
		RevokedAt: sql.NullInt64{},
		UserId:    sql.NullString{},
		TeamId:    sql.NullString{},
	}

	var target *model.NotePermissions
	if resolvedUserID != "" {
		permission.UserId = sql.NullString{String: resolvedUserID, Valid: true}
		target, err = l.svcCtx.NotePermissionsModel.FindOneByNoteIdUserId(l.ctx, in.NoteId, permission.UserId)
	} else {
		permission.TeamId = sql.NullString{String: in.TargetTeamId, Valid: true}
		target, err = l.svcCtx.NotePermissionsModel.FindOneByNoteIdTeamId(l.ctx, in.NoteId, permission.TeamId)
	}
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	if err == model.ErrNotFound {
		permission.PermissionId = strings.ReplaceAll(uuid.NewString(), "-", "")
		res, insertErr := l.svcCtx.NotePermissionsModel.Insert(l.ctx, permission)
		if insertErr != nil {
			return nil, insertErr
		}
		rows, rowsErr := res.RowsAffected()
		if rowsErr != nil {
			return nil, rowsErr
		}
		if rows != 1 {
			return nil, errors.New("failed to insert permission: no rows affected")
		}
	} else {
		target.GrantedBy = in.OperatorId
		target.Role = role
		target.Status = status
		target.GrantedAt = now
		target.RevokedAt = sql.NullInt64{}
		if updateErr := l.svcCtx.NotePermissionsModel.Update(l.ctx, target); updateErr != nil {
			return nil, updateErr
		}
		permission = target
	}
	payload := permissionPayload(permission.UserId.String, permission.TeamId.String, permission.Role)
	if err = appendCollaborationEvent(l.ctx, l.svcCtx, permission.NoteId, "permission_granted", in.OperatorId, payload, nil, nil, false); err != nil {
		return nil, err
	}

	return &syncnoterpc.PermissionResp{
		Success: true,
		Message: "Success",
		Permission: &syncnoterpc.PermissionInfo{
			PermissionId: permission.PermissionId,
			NoteId:       permission.NoteId,
			UserId:       permission.UserId.String,
			TeamId:       permission.TeamId.String,
			GrantedBy:    permission.GrantedBy,
			Role:         dbToRole(permission.Role),
			Status:       dbToStatus(permission.Status),
			GrantedAt:    permission.GrantedAt,
			RevokedAt:    permission.RevokedAt.Int64,
		},
	}, nil
}

func (l *GrantPermissionLogic) resolveUserIDByEmail(email string) (string, error) {
	type userRow struct {
		ID string `db:"id"`
	}

	var row userRow
	err := l.svcCtx.Conn.QueryRowCtx(l.ctx, &row, "select id from users where email = ? limit 1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("target user not found by email")
		}
		return "", err
	}
	if strings.TrimSpace(row.ID) == "" {
		return "", errors.New("target user not found by email")
	}
	return row.ID, nil
}

func roleToDB(role syncnoterpc.Role) (string, bool) {
	switch role {
	case syncnoterpc.Role_ROLE_OWNER:
		return "owner", true
	case syncnoterpc.Role_ROLE_ADMIN:
		return "admin", true
	case syncnoterpc.Role_ROLE_EDITOR:
		return "editor", true
	case syncnoterpc.Role_ROLE_VIEWER:
		return "viewer", true
	default:
		return "", false
	}
}

func dbToRole(role string) syncnoterpc.Role {
	switch strings.ToLower(role) {
	case "owner":
		return syncnoterpc.Role_ROLE_OWNER
	case "admin":
		return syncnoterpc.Role_ROLE_ADMIN
	case "editor":
		return syncnoterpc.Role_ROLE_EDITOR
	case "viewer":
		return syncnoterpc.Role_ROLE_VIEWER
	default:
		return syncnoterpc.Role_ROLE_UNSPECIFIED
	}
}

func dbToStatus(status string) syncnoterpc.PermissionStatus {
	switch strings.ToLower(status) {
	case "active":
		return syncnoterpc.PermissionStatus_PERMISSION_STATUS_ACTIVE
	case "revoked":
		return syncnoterpc.PermissionStatus_PERMISSION_STATUS_REVOKED
	case "pending":
		return syncnoterpc.PermissionStatus_PERMISSION_STATUS_PENDING
	default:
		return syncnoterpc.PermissionStatus_PERMISSION_STATUS_UNSPECIFIED
	}
}
