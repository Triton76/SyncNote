package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTeamMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取团队成员列表
func NewListTeamMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTeamMembersLogic {
	return &ListTeamMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTeamMembersLogic) ListTeamMembers(req *types.ListTeamMembersRequest) (resp *types.ListTeamMembersResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
