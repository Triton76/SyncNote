package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &syncnoterpc.RevokeNotePermissionResponse{}, nil
}
