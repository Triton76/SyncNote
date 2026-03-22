package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取笔记详情
func NewGetNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteLogic {
	return &GetNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoteLogic) GetNote(req *types.GetNoteRequest) (resp *types.GetNoteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
