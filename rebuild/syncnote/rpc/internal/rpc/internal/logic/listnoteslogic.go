package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotesLogic {
	return &ListNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNotesLogic) ListNotes(in *syncnoterpc.ListNotesRequest) (*syncnoterpc.ListNotesResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.ListNotesResponse{}, nil
}
