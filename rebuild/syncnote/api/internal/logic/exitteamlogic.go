package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExitTeamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 退出团队
func NewExitTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExitTeamLogic {
	return &ExitTeamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExitTeamLogic) ExitTeam(req *types.ExitTeamRequest) (resp *types.ExitTeamResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
