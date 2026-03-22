package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

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
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.ExitTeam(rpcCtx, &syncnoterpc.ExitTeamRequest{TeamId: req.TeamId})
	if err != nil {
		return nil, err
	}

	return &types.ExitTeamResponse{Success: rpcResp.GetSuccess()}, nil
}
