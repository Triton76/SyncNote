package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTeamMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTeamMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTeamMembersLogic {
	return &ListTeamMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTeamMembersLogic) ListTeamMembers(in *syncnoterpc.ListTeamMembersRequest) (*syncnoterpc.ListTeamMembersResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.ListTeamMembersResponse{}, nil
}
