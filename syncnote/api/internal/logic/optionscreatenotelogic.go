package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsCreateNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsCreateNoteLogic {
	return &OptionsCreateNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsCreateNoteLogic) OptionsCreateNote(req *types.UserNotesReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
