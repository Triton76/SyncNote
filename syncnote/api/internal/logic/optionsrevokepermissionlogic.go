package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsRevokePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsRevokePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsRevokePermissionLogic {
	return &OptionsRevokePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsRevokePermissionLogic) OptionsRevokePermission(req *types.UserNotesReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
