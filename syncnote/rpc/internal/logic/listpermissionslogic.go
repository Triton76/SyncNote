package logic

import (
	"context"

	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取笔记的所有权限列表
func (l *ListPermissionsLogic) ListPermissions(in *syncnoterpc.ListPermissionsReq) (*syncnoterpc.ListPermissionsResp, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.ListPermissionsResp{}, nil
}
