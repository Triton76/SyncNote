package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExitTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExitTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExitTeamLogic {
	return &ExitTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExitTeamLogic) ExitTeam(in *syncnoterpc.ExitTeamRequest) (*syncnoterpc.ExitTeamResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.ExitTeamResponse{}, nil
}
