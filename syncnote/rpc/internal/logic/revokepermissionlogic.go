package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &syncnoterpc.PermissionResp{}, nil
}
