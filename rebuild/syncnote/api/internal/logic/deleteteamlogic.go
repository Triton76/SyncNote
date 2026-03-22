package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTeamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除团队
func NewDeleteTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTeamLogic {
	return &DeleteTeamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTeamLogic) DeleteTeam(req *types.DeleteTeamRequest) (resp *types.DeleteTeamResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
