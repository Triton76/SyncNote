package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoteModel = (*customNoteModel)(nil)

type (
	// NoteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoteModel.
	NoteModel interface {
		noteModel
		withSession(session sqlx.Session) NoteModel
		UpdateWithVersion(ctx context.Context, note *Note, oldVersion int64) error
	}

	customNoteModel struct {
		*defaultNoteModel
	}
)

// NewNoteModel returns a model for the database table.
func NewNoteModel(conn sqlx.SqlConn) NoteModel {
	return &customNoteModel{
		defaultNoteModel: newNoteModel(conn),
	}
}

func (m *customNoteModel) withSession(session sqlx.Session) NoteModel {
	return NewNoteModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customNoteModel) UpdateWithVersion(ctx context.Context, note *Note, oldVersion int64) error {
	// 这个方法保证了乐观锁的并发安全
	query := fmt.Sprintf("update %s set title = ?, content = ?, version = version + 1 where note_id = ? and version = ? and deleted_at is null", m.table)

	res, err := m.conn.ExecCtx(ctx, query, note.Title, note.Content, note.NoteId, oldVersion)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrOptimisticLockFailed
	}
	return nil
}
