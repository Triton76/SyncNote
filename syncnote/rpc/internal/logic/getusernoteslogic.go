package logic

import (
	"context"
	"errors"

	"SyncNote/rpc/internal/svc"
	"SyncNote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNotesLogic {
	return &GetUserNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserNotesLogic) GetUserNotes(in *syncnoterpc.UserNotesReq) (*syncnoterpc.UserNotesResp, error) {
	if in.GetUserId() == "" {
		return nil, errors.New("userId is required")
	}

	notes, err := l.svcCtx.NoteStore.GetNotesByUserID(l.ctx, in.GetUserId())
	if err != nil {
		return nil, err
	}

	items := make([]*syncnoterpc.NoteSummary, 0, len(notes))
	for _, note := range notes {
		items = append(items, toNoteSummary(note))
	}

	return &syncnoterpc.UserNotesResp{Notes: items}, nil
}
