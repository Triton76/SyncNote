package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NotesModel = (*customNotesModel)(nil)

type (
	// NotesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotesModel.
	NotesModel interface {
		notesModel
		DumpUserNotes(ctx context.Context, userid string) ([]*Notes, error)
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

func (m *customNotesModel)DumpUserNotes(ctx context.Context, userid string) ([]*Notes, error) {
	//获取该用户所属所有笔记
	var resp []*Notes

	query := fmt.Sprintf("select %s from %s where user_id = ?", notesRows, m.table)
	
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userid)

	if err != nil {
		return nil, err
	}
	return resp, nil

}