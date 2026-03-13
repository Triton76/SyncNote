package logic

import (
	"context"
	"errors"

	"SyncNote/model"
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
	if in.GetUserId() == "" {
		return nil, errors.New("userId is required")
	}
	if in.GetTitle() == "" {
		return nil, errors.New("title is required")
	}

	created, err := l.svcCtx.NoteStore.CreateNote(l.ctx, &model.Note{
		UserID:  in.GetUserId(),
		Title:   in.GetTitle(),
		Content: in.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return toNoteResp(created), nil
}
