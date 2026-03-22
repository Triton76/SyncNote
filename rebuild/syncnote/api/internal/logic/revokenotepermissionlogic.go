package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeNotePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤销笔记权限
func NewRevokeNotePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeNotePermissionLogic {
	return &RevokeNotePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeNotePermissionLogic) RevokeNotePermission(req *types.RevokeNotePermissionRequest) (resp *types.RevokeNotePermissionResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
