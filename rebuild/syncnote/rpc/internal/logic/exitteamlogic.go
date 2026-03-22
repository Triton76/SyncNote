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

type ExitTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExitTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExitTeamLogic {
	return &ExitTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExitTeamLogic) ExitTeam(in *syncnoterpc.ExitTeamRequest) (*syncnoterpc.ExitTeamResponse, error) {
	// 当前用户离开指定团队，删除 team_members 记录。
	if in.GetTeamId() == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}

	userId, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	member, err := l.svcCtx.TeamMembersModel.FindOneByTeamIdUserId(l.ctx, in.GetTeamId(), userId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "member not found in team")
		}
		return nil, err
	}

	if err := l.svcCtx.TeamMembersModel.Delete(l.ctx, member.Id); err != nil {
		return nil, err
	}

	return &syncnoterpc.ExitTeamResponse{Success: true}, nil
}
