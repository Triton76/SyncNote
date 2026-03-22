package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

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
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.UpdateTeam(rpcCtx, &syncnoterpc.UpdateTeamRequest{
		TeamId:      req.TeamId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateTeamResponse{Team: toAPITeam(rpcResp.GetTeam())}, nil
}
