package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantNotePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 授予笔记权限
func NewGrantNotePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantNotePermissionLogic {
	return &GrantNotePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GrantNotePermissionLogic) GrantNotePermission(req *types.GrantNotePermissionRequest) (resp *types.GrantNotePermissionResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
