package logic

import (
	"context"
	"time"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JoinTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinTeamLogic {
	return &JoinTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinTeamLogic) JoinTeam(in *syncnoterpc.JoinTeamRequest) (*syncnoterpc.JoinTeamResponse, error) {
	// 当前用户加入指定团队，添加 team_members 记录。
	if in.GetTeamId() == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}

	userId, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	// 检查团队存在
	_, err = l.svcCtx.TeamModel.FindOne(l.ctx, in.GetTeamId())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "team not found")
		}
		return nil, err
	}

	// 检查是否已经在团队中
	_, err = l.svcCtx.TeamMembersModel.FindOneByTeamIdUserId(l.ctx, in.GetTeamId(), userId)
	if err != model.ErrNotFound && err != nil {
		return nil, err
	}
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user already in team")
	}

	// 添加成员
	member := &model.TeamMembers{
		Id:       uuid.NewString(),
		TeamId:   in.GetTeamId(),
		UserId:   userId,
		JoinedAt: time.Now(),
	}
	_, err = l.svcCtx.TeamMembersModel.Insert(l.ctx, member)
	if err != nil {
		return nil, err
	}

	return &syncnoterpc.JoinTeamResponse{Success: true}, nil
}
