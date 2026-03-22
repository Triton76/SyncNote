package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinTeamLogic {
	return &JoinTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinTeamLogic) JoinTeam(in *syncnoterpc.JoinTeamRequest) (*syncnoterpc.JoinTeamResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.JoinTeamResponse{}, nil
}
