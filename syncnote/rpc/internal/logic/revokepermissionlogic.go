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

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokePermissionLogic {
	return &RevokePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤销权限
func (l *RevokePermissionLogic) RevokePermission(in *syncnoterpc.RevokePermissionReq) (*syncnoterpc.PermissionResp, error) {
	if in == nil {
		return nil, errors.New("request is nil")
	}
	if in.NoteId == "" || in.OperatorId == "" {
		return nil, errors.New("noteId and operatorId are required")
	}
	if (in.TargetUserId == "" && in.TargetTeamId == "") || (in.TargetUserId != "" && in.TargetTeamId != "") {
		return nil, errors.New("targetUserId and targetTeamId must be exactly one")
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

	var target *model.NotePermissions
	if in.TargetUserId != "" {
		target, err = l.svcCtx.NotePermissionsModel.FindOneByNoteIdUserId(l.ctx, in.NoteId, sql.NullString{String: in.TargetUserId, Valid: true})
	} else {
		target, err = l.svcCtx.NotePermissionsModel.FindOneByNoteIdTeamId(l.ctx, in.NoteId, sql.NullString{String: in.TargetTeamId, Valid: true})
	}
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	if strings.ToLower(target.Role) == "owner" {
		return nil, errors.New("cannot revoke owner role")
	}

	now := time.Now().UnixMilli()
	target.Status = "revoked"
	target.RevokedAt = sql.NullInt64{Int64: now, Valid: true}
	if err = l.svcCtx.NotePermissionsModel.Update(l.ctx, target); err != nil {
		return nil, err
	}
	payload := permissionPayload(target.UserId.String, target.TeamId.String, target.Role)
	if err = appendCollaborationEvent(l.ctx, l.svcCtx, target.NoteId, "permission_revoked", in.OperatorId, payload, nil, nil, false); err != nil {
		return nil, err
	}

	return &syncnoterpc.PermissionResp{
		Success: true,
		Message: "Success",
		Permission: &syncnoterpc.PermissionInfo{
			PermissionId: target.PermissionId,
			NoteId:       target.NoteId,
			UserId:       target.UserId.String,
			TeamId:       target.TeamId.String,
			GrantedBy:    target.GrantedBy,
			Role:         dbToRole(target.Role),
			Status:       dbToStatus(target.Status),
			GrantedAt:    target.GrantedAt,
			RevokedAt:    target.RevokedAt.Int64,
		},
	}, nil
}
