package logic

import (
	"context"
	"errors"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"SyncNote/syncnote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNoteEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteEventsLogic {
	return &GetNoteEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoteEventsLogic) GetNoteEvents(req *types.GetNoteEventsReq) (resp *types.GetNoteEventsResp, err error) {
	if _, err = currentUserIDFromCtx(l.ctx); err != nil {
		return nil, err
	}
	if req.NoteId == "" {
		return nil, errors.New("noteId is required")
	}

	rpcResp, err := l.svcCtx.SyncNoteRpc.GetNoteEvents(l.ctx, &syncnoterpcclient.GetNoteEventsReq{
		NoteId:   req.NoteId,
		StartSeq: req.StartSeq,
		Limit:    req.Limit,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return &types.GetNoteEventsResp{}, nil
	}

	events := make([]types.CollaborationEvent, 0, len(rpcResp.Events))
	for _, e := range rpcResp.Events {
		events = append(events, types.CollaborationEvent{
			EventId:      e.EventId,
			NoteId:       e.NoteId,
			EventSeq:     e.EventSeq,
			EventType:    eventTypeToString(e.EventType),
			OperatorId:   e.OperatorId,
			OperatorName: e.OperatorName,
			Payload:      e.Payload,
			NoteVersion:  e.NoteVersion,
			IsConflict:   e.IsConflict,
			CreatedAt:    e.CreatedAt,
		})
	}

	return &types.GetNoteEventsResp{Events: events, HasMore: rpcResp.HasMore}, nil
}
