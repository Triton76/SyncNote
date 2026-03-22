package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTeamLogic {
	return &CreateTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTeamLogic) CreateTeam(in *syncnoterpc.CreateTeamRequest) (*syncnoterpc.CreateTeamResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.CreateTeamResponse{}, nil
}
