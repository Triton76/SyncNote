package logic

import (
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"
	"context"
	"errors"

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

	note, err := l.svcCtx.NotesModel.FindOne(l.ctx, in.GetNoteId())
	if err != nil {
		return nil, err
	}
	return &syncnoterpc.NoteResp{
		NoteId: note.NoteId,
		UserId: note.UserId,
		Title: note.Title,
		Content: note.Content,
		Version: int64(note.Version),
		LastModified: note.LastModified,
	}, nil
}
