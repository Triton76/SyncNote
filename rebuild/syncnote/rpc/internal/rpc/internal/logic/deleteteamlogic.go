package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTeamLogic {
	return &DeleteTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTeamLogic) DeleteTeam(in *syncnoterpc.DeleteTeamRequest) (*syncnoterpc.DeleteTeamResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.DeleteTeamResponse{}, nil
}
