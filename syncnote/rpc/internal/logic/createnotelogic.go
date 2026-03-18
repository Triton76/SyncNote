package logic

import (
	"SyncNote/syncnote/rpc/internal/model"
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	if in.GetUserId() == "" {
		return nil, errors.New("userId is required")
	}
	if in.GetTitle() == "" {
		return nil, errors.New("title is required")
	}
	newNoteId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	note := &model.Notes{
		NoteId:       newNoteId.String(),
		UserId:       in.UserId,
		Title:        in.Title,
		Content:      in.Content,
		Version:      1,
		LastModified: time.Now().Unix(),
	}
	res, err := l.svcCtx.NotesModel.Insert(l.ctx, note)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("save note failed")
	}
	return &syncnoterpc.NoteResp{
		NoteId:       note.NoteId,
		UserId:       note.UserId,
		Title:        note.Title,
		Content:      note.Content,
		Version:      int64(note.Version),
		LastModified: note.LastModified,
	}, nil
}
