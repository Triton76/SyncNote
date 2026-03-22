package logic

import (
	"context"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RevokeNotePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeNotePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeNotePermissionLogic {
	return &RevokeNotePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RevokeNotePermissionLogic) RevokeNotePermission(in *syncnoterpc.RevokeNotePermissionRequest) (*syncnoterpc.RevokeNotePermissionResponse, error) {
	// 撤销权限：仅 owner 或 admin 权限用户可操作。
	if in.GetNoteId() == "" || in.GetTargetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "note_id and target_user_id are required")
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

	// 鉴权：owner 或 admin 权限
	if note.OwnerId != operatorID {
		ok, checkErr := l.canRevokeByAdmin(in.GetNoteId(), operatorID)
		if checkErr != nil {
			return nil, checkErr
		}
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	// 删除权限
	perm, err := l.svcCtx.NoteUserPermissionModel.FindOneByNoteIdUserId(l.ctx, in.GetNoteId(), in.GetTargetUserId())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		return nil, err
	}

	err = l.svcCtx.NoteUserPermissionModel.Delete(l.ctx, perm.PermissionId)
	if err != nil {
		return nil, err
	}

	return &syncnoterpc.RevokeNotePermissionResponse{Success: true}, nil
}

func (l *RevokeNotePermissionLogic) canRevokeByAdmin(noteId, operatorID string) (bool, error) {
	userPerm, err := l.svcCtx.NoteUserPermissionModel.FindOneByNoteIdUserId(l.ctx, noteId, operatorID)
	if err != nil && err != model.ErrNotFound {
		return false, err
	}
	if err == nil && userPerm.PermissionLevel == model.PermissionLevelAdmin {
		return true, nil
	}

	return l.svcCtx.NoteTeamPermissionModel.ExistsTeamPermissionLevel(l.ctx, noteId, operatorID, model.PermissionLevelAdmin)
}
