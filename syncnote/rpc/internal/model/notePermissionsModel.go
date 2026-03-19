package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NotePermissionsModel = (*customNotePermissionsModel)(nil)

type (
	// NotePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotePermissionsModel.
	NotePermissionsModel interface {
		notePermissionsModel
		ListByNoteId(ctx context.Context, noteId string) ([]*NotePermissions, error)
	}

	customNotePermissionsModel struct {
		*defaultNotePermissionsModel
	}
)

// NewNotePermissionsModel returns a model for the database table.
func NewNotePermissionsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NotePermissionsModel {
	return &customNotePermissionsModel{
		defaultNotePermissionsModel: newNotePermissionsModel(conn, c, opts...),
	}
}

func (m *customNotePermissionsModel) ListByNoteId(ctx context.Context, noteId string) ([]*NotePermissions, error) {
	var resp []*NotePermissions
	query := fmt.Sprintf("select %s from %s where `note_id` = ?", notePermissionsRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &resp, query, noteId); err != nil {
		return nil, err
	}
	return resp, nil
}
