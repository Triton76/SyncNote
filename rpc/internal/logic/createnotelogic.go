package logic

import (
	"context"

	"SyncNote/rpc/internal/svc"
	"SyncNote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoteLogic {
	return &CreateNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNoteLogic) CreateNote(in *syncnoterpc.CreateNoteReq) (*syncnoterpc.NoteResp, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.NoteResp{}, nil
}
