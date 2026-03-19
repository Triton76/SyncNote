package logic

import (
	"context"
	"errors"
	"strings"

	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteEventsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoteEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteEventsLogic {
	return &GetNoteEventsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Collaboration History (对应 collaboration_events 表) ---
func (l *GetNoteEventsLogic) GetNoteEvents(in *syncnoterpc.GetNoteEventsReq) (*syncnoterpc.GetNoteEventsResp, error) {
	if in == nil {
		return nil, errors.New("request is nil")
	}
	if in.NoteId == "" {
		return nil, errors.New("noteId is required")
	}

	startSeq := in.StartSeq
	if startSeq < 0 {
		startSeq = 0
	}

	limit := int64(in.Limit)
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	rows, err := l.svcCtx.CollabEventsModel.ListByNoteIDStartSeq(l.ctx, in.NoteId, startSeq, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := false
	if int64(len(rows)) > limit {
		hasMore = true
		rows = rows[:limit]
	}

	events := make([]*syncnoterpc.CollaborationEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, &syncnoterpc.CollaborationEvent{
			EventId:      row.EventId,
			NoteId:       row.NoteId,
			EventSeq:     row.EventSeq,
			EventType:    dbToEventType(row.EventType),
			OperatorId:   row.OperatorId,
			OperatorName: row.OperatorName.String,
			Payload:      row.Payload.String,
			NoteVersion:  row.NoteVersion.Int64,
			IsConflict:   row.IsConflict != 0,
			CreatedAt:    row.CreatedAt,
		})
	}

	return &syncnoterpc.GetNoteEventsResp{
		Events:  events,
		HasMore: hasMore,
	}, nil
}

func dbToEventType(eventType string) syncnoterpc.EventType {
	switch strings.ToLower(eventType) {
	case "note_created":
		return syncnoterpc.EventType_EVENT_TYPE_NOTE_CREATED
	case "note_updated":
		return syncnoterpc.EventType_EVENT_TYPE_NOTE_UPDATED
	case "note_deleted":
		return syncnoterpc.EventType_EVENT_TYPE_NOTE_DELETED
	case "permission_granted":
		return syncnoterpc.EventType_EVENT_TYPE_PERMISSION_GRANTED
	case "permission_revoked":
		return syncnoterpc.EventType_EVENT_TYPE_PERMISSION_REVOKED
	case "conflict_detected":
		return syncnoterpc.EventType_EVENT_TYPE_CONFLICT_DETECTED
	case "view_started":
		return syncnoterpc.EventType_EVENT_TYPE_VIEW_STARTED
	case "view_ended":
		return syncnoterpc.EventType_EVENT_TYPE_VIEW_ENDED
	default:
		return syncnoterpc.EventType_EVENT_TYPE_UNSPECIFIED
	}
}
