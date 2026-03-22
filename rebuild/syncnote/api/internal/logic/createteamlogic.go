package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTeamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建团队
func NewCreateTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTeamLogic {
	return &CreateTeamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTeamLogic) CreateTeam(req *types.CreateTeamRequest) (resp *types.CreateTeamResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
