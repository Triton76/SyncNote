package logic

import (
	"context"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTeamLogic {
	return &GetTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTeamLogic) GetTeam(in *syncnoterpc.GetTeamRequest) (*syncnoterpc.GetTeamResponse, error) {
	// 每个用户都可以查看团队信息，暂时不考虑权限问题

	if in.GetTeamId() == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}

	_, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	team, err := l.svcCtx.TeamModel.FindOne(l.ctx, in.GetTeamId())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "team not found")
		}
		return nil, err
	}

	return &syncnoterpc.GetTeamResponse{Team: &syncnoterpc.Team{
		TeamId:      team.TeamId,
		Name:        team.Name,
		Description: team.Description.String,
		OwnerId:     team.OwnerId,
	}}, nil
}
