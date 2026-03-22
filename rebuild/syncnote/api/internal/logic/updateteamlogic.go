package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTeamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新团队
func NewUpdateTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTeamLogic {
	return &UpdateTeamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTeamLogic) UpdateTeam(req *types.UpdateTeamRequest) (resp *types.UpdateTeamResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
