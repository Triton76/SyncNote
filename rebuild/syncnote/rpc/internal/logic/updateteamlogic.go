package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTeamLogic {
	return &UpdateTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTeamLogic) UpdateTeam(in *syncnoterpc.UpdateTeamRequest) (*syncnoterpc.UpdateTeamResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.UpdateTeamResponse{}, nil
}
