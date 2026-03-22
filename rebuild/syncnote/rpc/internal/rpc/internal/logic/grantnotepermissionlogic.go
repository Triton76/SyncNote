package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &syncnoterpc.GrantNotePermissionResponse{}, nil
}
