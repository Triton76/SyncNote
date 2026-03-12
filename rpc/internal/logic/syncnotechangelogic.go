package logic

import (
	"context"

	"SyncNote/rpc/internal/svc"
	"SyncNote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncNoteChangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncNoteChangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncNoteChangeLogic {
	return &SyncNoteChangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncNoteChangeLogic) SyncNoteChange(in *syncnoterpc.NoteChangeReq) (*syncnoterpc.SyncResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.SyncResponse{}, nil
}
