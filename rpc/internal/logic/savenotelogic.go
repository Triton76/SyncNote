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

func (l *SaveNoteLogic) SaveNote(in *syncnoterpc.SaveNoteReq) (*syncnoterpc.SaveNoteResp, error) {
	if in.GetNoteId() == "" {
		return invalidParamSaveResp("noteId is required"), nil
	}
	if in.GetUserId() == "" {
		return invalidParamSaveResp("userId is required"), nil
	}
	if in.GetExpectedVersion() <= 0 {
		return invalidParamSaveResp("expectedVersion must be greater than 0"), nil
	}

	result, err := l.svcCtx.NoteStore.SaveNote(
		l.ctx,
		in.GetNoteId(),
		in.GetUserId(),
		in.GetContent(),
		in.GetExpectedVersion(),
	)
	if err != nil {
		return nil, err
	}

	return toSaveNoteResp(result), nil
}
