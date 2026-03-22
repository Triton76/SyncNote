package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取笔记列表
func NewListNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotesLogic {
	return &ListNotesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNotesLogic) ListNotes(req *types.ListNotesRequest) (resp *types.ListNotesResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
