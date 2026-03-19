package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NotePermissionsModel = (*customNotePermissionsModel)(nil)

type (
	// NotePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotePermissionsModel.
	NotePermissionsModel interface {
		notePermissionsModel
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
