package logic

import (
	"context"
	"errors"

	"SyncNote/rpc/internal/svc"
	"SyncNote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteLogic {
	return &GetNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoteLogic) GetNote(in *syncnoterpc.NoteReq) (*syncnoterpc.NoteResp, error) {
	if in.GetNoteId() == "" {
		return nil, errors.New("noteId is required")
	}

	note, err := l.svcCtx.NoteStore.GetNoteByID(l.ctx, in.GetNoteId())
	if err != nil {
		return nil, err
	}

	return toNoteResp(note), nil
}
