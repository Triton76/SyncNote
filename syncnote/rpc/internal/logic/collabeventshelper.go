package logic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"SyncNote/syncnote/rpc/internal/model"
	"SyncNote/syncnote/rpc/internal/svc"

	"github.com/google/uuid"
)

func appendCollaborationEvent(ctx context.Context, svcCtx *svc.ServiceContext, noteID, eventType, operatorID, payload string, noteVersion, expectedVersion *int64, isConflict bool) error {
	e := &model.CollaborationEvents{
		EventId:         strings.ReplaceAll(uuid.NewString(), "-", ""),
		NoteId:          noteID,
		EventSeq:        time.Now().UnixNano(),
		EventType:       eventType,
		OperatorId:      operatorID,
		OperatorName:    sql.NullString{},
		Payload:         sql.NullString{},
		NoteVersion:     sql.NullInt64{},
		ExpectedVersion: sql.NullInt64{},
		IsConflict:      0,
		RelatedEventId:  sql.NullString{},
		CreatedAt:       time.Now().UnixMilli(),
	}
	if payload != "" {
		e.Payload = sql.NullString{String: payload, Valid: true}
	}
	if noteVersion != nil {
		e.NoteVersion = sql.NullInt64{Int64: *noteVersion, Valid: true}
	}
	if expectedVersion != nil {
		e.ExpectedVersion = sql.NullInt64{Int64: *expectedVersion, Valid: true}
	}
	if isConflict {
		e.IsConflict = 1
	}
	_, err := svcCtx.CollabEventsModel.InsertEvent(ctx, e)
	return err
}

func permissionPayload(targetUserID, targetTeamID, role string) string {
	return fmt.Sprintf("{\"targetUserId\":\"%s\",\"targetTeamId\":\"%s\",\"role\":\"%s\"}", targetUserID, targetTeamID, role)
}
