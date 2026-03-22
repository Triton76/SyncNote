package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTeamLogic {
	return &GetTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTeamLogic) GetTeam(in *syncnoterpc.GetTeamRequest) (*syncnoterpc.GetTeamResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.GetTeamResponse{}, nil
}
