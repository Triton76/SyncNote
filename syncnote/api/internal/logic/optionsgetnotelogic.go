package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsGetNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsGetNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsGetNoteLogic {
	return &OptionsGetNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsGetNoteLogic) OptionsGetNote(req *types.NoteReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
