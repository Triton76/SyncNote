package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CollaborationEventsModel = (*customCollaborationEventsModel)(nil)

type (
	// CollaborationEventsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCollaborationEventsModel.
	CollaborationEventsModel interface {
		collaborationEventsModel
		ListByNoteIDStartSeq(ctx context.Context, noteID string, startSeq int64, limit int64) ([]*CollaborationEvents, error)
		InsertEvent(ctx context.Context, data *CollaborationEvents) (sql.Result, error)
	}

	customCollaborationEventsModel struct {
		*defaultCollaborationEventsModel
	}
)

// NewCollaborationEventsModel returns a model for the database table.
func NewCollaborationEventsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CollaborationEventsModel {
	return &customCollaborationEventsModel{
		defaultCollaborationEventsModel: newCollaborationEventsModel(conn, c, opts...),
	}
}

func (m *customCollaborationEventsModel) ListByNoteIDStartSeq(ctx context.Context, noteID string, startSeq int64, limit int64) ([]*CollaborationEvents, error) {
	var resp []*CollaborationEvents
	query := fmt.Sprintf("select %s from %s where `note_id` = ? and `event_seq` >= ? order by `event_seq` asc limit ?", collaborationEventsRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &resp, query, noteID, startSeq, limit); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customCollaborationEventsModel) InsertEvent(ctx context.Context, data *CollaborationEvents) (sql.Result, error) {
	collaborationEventsEventIdKey := fmt.Sprintf("%s%v", cacheCollaborationEventsEventIdPrefix, data.EventId)
	collaborationEventsNoteIdEventSeqKey := fmt.Sprintf("%s%v:%v", cacheCollaborationEventsNoteIdEventSeqPrefix, data.NoteId, data.EventSeq)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("insert into %s (`event_id`,`note_id`,`event_seq`,`event_type`,`operator_id`,`operator_name`,`payload`,`note_version`,`expected_version`,`is_conflict`,`related_event_id`,`created_at`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
		return conn.ExecCtx(ctx, query,
			data.EventId,
			data.NoteId,
			data.EventSeq,
			data.EventType,
			data.OperatorId,
			data.OperatorName,
			data.Payload,
			data.NoteVersion,
			data.ExpectedVersion,
			data.IsConflict,
			data.RelatedEventId,
			data.CreatedAt,
		)
	}, collaborationEventsEventIdKey, collaborationEventsNoteIdEventSeqKey)
	return ret, err
}
