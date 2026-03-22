package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTeamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取团队详情
func NewGetTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTeamLogic {
	return &GetTeamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTeamLogic) GetTeam(req *types.GetTeamRequest) (resp *types.GetTeamResponse, err error) {
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.GetTeam(rpcCtx, &syncnoterpc.GetTeamRequest{TeamId: req.TeamId})
	if err != nil {
		return nil, err
	}

	return &types.GetTeamResponse{Team: toAPITeam(rpcResp.GetTeam())}, nil
}
