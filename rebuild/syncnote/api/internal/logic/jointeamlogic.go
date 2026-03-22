package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

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
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.JoinTeam(rpcCtx, &syncnoterpc.JoinTeamRequest{TeamId: req.TeamId})
	if err != nil {
		return nil, err
	}

	return &types.JoinTeamResponse{Success: rpcResp.GetSuccess()}, nil
}
