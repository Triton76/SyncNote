package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NotesModel = (*customNotesModel)(nil)

type (
	// NotesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotesModel.
	NotesModel interface {
		notesModel
	}

	customNotesModel struct {
		*defaultNotesModel
	}
)

// NewNotesModel returns a model for the database table.
func NewNotesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NotesModel {
	return &customNotesModel{
		defaultNotesModel: newNotesModel(conn, c, opts...),
	}
}
