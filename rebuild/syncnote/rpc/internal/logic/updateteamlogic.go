package logic

import (
	"context"
	"database/sql"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTeamLogic {
	return &UpdateTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTeamLogic) UpdateTeam(in *syncnoterpc.UpdateTeamRequest) (*syncnoterpc.UpdateTeamResponse, error) {
	// 更新团队：仅 owner 可更新。
	if in.GetTeamId() == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}

	operatorID, err := middleware.GetUserFromContext(l.ctx)
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

	// 鉴权：仅 owner 可修改
	if team.OwnerId != operatorID {
		return nil, status.Error(codes.PermissionDenied, "only team owner can update team")
	}

	// 更新团队
	team.Name = in.GetName()
	team.Description = sql.NullString{String: in.GetDescription(), Valid: in.GetDescription() != ""}

	err = l.svcCtx.TeamModel.Update(l.ctx, team)
	if err != nil {
		return nil, err
	}

	return &syncnoterpc.UpdateTeamResponse{Team: &syncnoterpc.Team{
		TeamId:      team.TeamId,
		Name:        team.Name,
		Description: team.Description.String,
		OwnerId:     team.OwnerId,
	}}, nil
}
