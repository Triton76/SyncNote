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

type DeleteTeamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTeamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTeamLogic {
	return &DeleteTeamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTeamLogic) DeleteTeam(in *syncnoterpc.DeleteTeamRequest) (*syncnoterpc.DeleteTeamResponse, error) {
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

	if team.OwnerId != operatorID {
		return nil, status.Error(codes.PermissionDenied, "only team owner can delete team")
	}

	if err := l.svcCtx.TeamModel.Delete(l.ctx, in.GetTeamId()); err != nil {
		return nil, err
	}

	return &syncnoterpc.DeleteTeamResponse{Success: true}, nil
}
