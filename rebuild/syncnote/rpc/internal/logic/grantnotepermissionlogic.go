package logic

import (
	"context"
	"database/sql"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrantNotePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGrantNotePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantNotePermissionLogic {
	return &GrantNotePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GrantNotePermissionLogic) GrantNotePermission(in *syncnoterpc.GrantNotePermissionRequest) (*syncnoterpc.GrantNotePermissionResponse, error) {
	// 权限的起点为Owner。
	// 对于团队，只有权限为admin的团队才能决定其他用户对这个笔记的权限。
	// 对于用户，也只有权限为admin对用户才能决定其他用户的权限。
	// Owner拥有最高权限。
	/*  伪代码块
	operator_Id := context
	判断op是否有权利: {
		1. note.owner == operator? if not then
		2. 查 note_user_permission 是否存在(note_id, user_id=operator_id, permission_level='admin') if not then
		3. 查 operator 所在团队里，是否存在某team_id在note_team_permission上满足(note_id, team_id, permission_level='admin') if not then
		return 权限不足
	}
	*/
	if in.GetNoteId() == "" || in.GetTargetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "note_id and target_user_id are required")
	}

	level, err := permissionLevelToDB(in.GetLevel())
	if err != nil {
		return nil, err
	}

	operatorID, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	note, err := l.svcCtx.NoteModel.FindOne(l.ctx, in.GetNoteId())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "note not found")
		}
		return nil, err
	}

	if !canGrantByOwner(note.OwnerId, operatorID) {
		ok, checkErr := l.canGrantByAdmin(in.GetNoteId(), operatorID)
		if checkErr != nil {
			return nil, checkErr
		}
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	//鉴权通过准备授予权限，接下来看授予权限的对象存不存在。
	if _, err := l.svcCtx.UserModel.FindOne(l.ctx, in.GetTargetUserId()); err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "target user not found")
		}
		return nil, err
	}

	grantedBy := sql.NullString{String: operatorID, Valid: true}
	//接下来查有没有已经存在的权限记录
	existing, err := l.svcCtx.NoteUserPermissionModel.FindOneByNoteIdUserId(l.ctx, in.GetNoteId(), in.GetTargetUserId())
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	permission := &syncnoterpc.UserPermission{
		NoteId:    in.GetNoteId(),
		UserId:    in.GetTargetUserId(),
		Level:     in.GetLevel(),
		GrantedBy: operatorID,
	}

	//没查到对应权限记录，进行下面的if
	if err == model.ErrNotFound {
		record := &model.NoteUserPermission{
			PermissionId:    uuid.NewString(),
			NoteId:          in.GetNoteId(),
			UserId:          in.GetTargetUserId(),
			PermissionLevel: level,
			GrantedBy:       grantedBy,
		}
		if _, insertErr := l.svcCtx.NoteUserPermissionModel.Insert(l.ctx, record); insertErr != nil {
			return nil, insertErr
		}
		permission.PermissionId = record.PermissionId
	} else {
		//查到了就update
		existing.PermissionLevel = level
		existing.GrantedBy = grantedBy
		if updateErr := l.svcCtx.NoteUserPermissionModel.Update(l.ctx, existing); updateErr != nil {
			return nil, updateErr
		}
		permission.PermissionId = existing.PermissionId
	}

	return &syncnoterpc.GrantNotePermissionResponse{Permission: permission}, nil
}

func (l *GrantNotePermissionLogic) canGrantByAdmin(noteId, operatorID string) (bool, error) {
	userPerm, err := l.svcCtx.NoteUserPermissionModel.FindOneByNoteIdUserId(l.ctx, noteId, operatorID)
	if err != nil && err != model.ErrNotFound {
		return false, err
	}
	if err == nil && userPerm.PermissionLevel == "admin" {
		return true, nil
	}

	return l.svcCtx.NoteTeamPermissionModel.ExistsTeamAdminForUser(l.ctx, noteId, operatorID)
}

func canGrantByOwner(ownerID, operatorID string) bool {
	return ownerID == operatorID
}

func permissionLevelToDB(level syncnoterpc.PermissionLevel) (string, error) {
	switch level {
	case syncnoterpc.PermissionLevel_PERMISSION_LEVEL_READ:
		return "read", nil
	case syncnoterpc.PermissionLevel_PERMISSION_LEVEL_WRITE:
		return "write", nil
	case syncnoterpc.PermissionLevel_PERMISSION_LEVEL_ADMIN:
		return "admin", nil
	default:
		return "", status.Error(codes.InvalidArgument, "invalid permission level")
	}
}
