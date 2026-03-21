package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ NoteUserPermissionModel = (*customNoteUserPermissionModel)(nil)

type (
	// NoteUserPermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoteUserPermissionModel.
	NoteUserPermissionModel interface {
		noteUserPermissionModel
		withSession(session sqlx.Session) NoteUserPermissionModel
	}

	customNoteUserPermissionModel struct {
		*defaultNoteUserPermissionModel
	}
)

// NewNoteUserPermissionModel returns a model for the database table.
func NewNoteUserPermissionModel(conn sqlx.SqlConn) NoteUserPermissionModel {
	return &customNoteUserPermissionModel{
		defaultNoteUserPermissionModel: newNoteUserPermissionModel(conn),
	}
}

func (m *customNoteUserPermissionModel) withSession(session sqlx.Session) NoteUserPermissionModel {
	return NewNoteUserPermissionModel(sqlx.NewSqlConnFromSession(session))
}
