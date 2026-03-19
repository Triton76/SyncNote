package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsGetNoteEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsGetNoteEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsGetNoteEventsLogic {
	return &OptionsGetNoteEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsGetNoteEventsLogic) OptionsGetNoteEvents(req *types.GetNoteEventsReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
