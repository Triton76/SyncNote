package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"SyncNote/syncnote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyTeamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyTeamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyTeamsLogic {
	return &GetMyTeamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyTeamsLogic) GetMyTeams(req *types.MyTeamsReq) (resp *types.MyTeamsResp, err error) {
	userID, err := currentUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.SyncNoteRpc.ListMyTeams(l.ctx, &syncnoterpcclient.ListMyTeamsReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return &types.MyTeamsResp{}, nil
	}

	out := make([]types.TeamInfo, 0, len(rpcResp.Teams))
	for _, t := range rpcResp.Teams {
		out = append(out, types.TeamInfo{
			TeamId:   t.TeamId,
			TeamName: t.TeamName,
			Role:     t.Role,
			Status:   t.Status,
			JoinedAt: t.JoinedAt,
		})
	}

	return &types.MyTeamsResp{Teams: out}, nil
}
