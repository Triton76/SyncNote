package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinTeamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 加入团队
func NewJoinTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinTeamLogic {
	return &JoinTeamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinTeamLogic) JoinTeam(req *types.JoinTeamRequest) (resp *types.JoinTeamResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
