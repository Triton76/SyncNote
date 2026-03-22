package logic

import (
	"context"
	"database/sql"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTeamLogic {
	return &CreateTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTeamLogic) CreateTeam(in *syncnoterpc.CreateTeamRequest) (*syncnoterpc.CreateTeamResponse, error) {
	// 创建团队，插入数据库，设置owner为当前用户。
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	userId, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	description := sql.NullString{String: in.Description, Valid: in.Description != ""}
	team := &model.Team{
		TeamId:      uuid.NewString(),
		Name:        in.Name,
		Description: description,
		OwnerId:     userId,
	}
	_, err = l.svcCtx.TeamModel.Insert(l.ctx, team)
	if err != nil {
		return nil, err
	}
	return &syncnoterpc.CreateTeamResponse{Team: &syncnoterpc.Team{
		TeamId:      team.TeamId,
		Name:        team.Name,
		Description: team.Description.String,
		OwnerId:     team.OwnerId,
	}}, nil
}
