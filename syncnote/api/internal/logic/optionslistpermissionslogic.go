package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsListPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsListPermissionsLogic {
	return &OptionsListPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsListPermissionsLogic) OptionsListPermissions(req *types.ListPermissionsReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
