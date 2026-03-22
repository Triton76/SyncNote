package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取笔记权限列表
func NewListNotePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotePermissionsLogic {
	return &ListNotePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNotePermissionsLogic) ListNotePermissions(req *types.ListNotePermissionsRequest) (resp *types.ListNotePermissionsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
