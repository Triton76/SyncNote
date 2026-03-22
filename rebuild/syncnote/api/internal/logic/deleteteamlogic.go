package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

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
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.DeleteTeam(rpcCtx, &syncnoterpc.DeleteTeamRequest{TeamId: req.TeamId})
	if err != nil {
		return nil, err
	}

	return &types.DeleteTeamResponse{Success: rpcResp.GetSuccess()}, nil
}
