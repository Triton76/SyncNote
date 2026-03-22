package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoteLogic {
	return &DeleteNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNoteLogic) DeleteNote(in *syncnoterpc.DeleteNoteRequest) (*syncnoterpc.DeleteNoteResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.DeleteNoteResponse{}, nil
}
