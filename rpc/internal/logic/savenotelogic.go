package logic

import (
	"context"

	"SyncNote/rpc/internal/svc"
	"SyncNote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoteLogic {
	return &SaveNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveNoteLogic) SaveNote(in *syncnoterpc.NoteReq) (*syncnoterpc.NoteResp, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.NoteResp{}, nil
}
