package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsGetMyTeamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsGetMyTeamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsGetMyTeamsLogic {
	return &OptionsGetMyTeamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsGetMyTeamsLogic) OptionsGetMyTeams(req *types.MyTeamsReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
