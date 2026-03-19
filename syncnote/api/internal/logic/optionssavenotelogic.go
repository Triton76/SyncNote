package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsSaveNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsSaveNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsSaveNoteLogic {
	return &OptionsSaveNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsSaveNoteLogic) OptionsSaveNote(req *types.UserNotesReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
