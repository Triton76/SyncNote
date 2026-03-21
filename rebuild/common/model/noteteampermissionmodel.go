package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ NoteTeamPermissionModel = (*customNoteTeamPermissionModel)(nil)

type (
	// NoteTeamPermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoteTeamPermissionModel.
	NoteTeamPermissionModel interface {
		noteTeamPermissionModel
		withSession(session sqlx.Session) NoteTeamPermissionModel
	}

	customNoteTeamPermissionModel struct {
		*defaultNoteTeamPermissionModel
	}
)

// NewNoteTeamPermissionModel returns a model for the database table.
func NewNoteTeamPermissionModel(conn sqlx.SqlConn) NoteTeamPermissionModel {
	return &customNoteTeamPermissionModel{
		defaultNoteTeamPermissionModel: newNoteTeamPermissionModel(conn),
	}
}

func (m *customNoteTeamPermissionModel) withSession(session sqlx.Session) NoteTeamPermissionModel {
	return NewNoteTeamPermissionModel(sqlx.NewSqlConnFromSession(session))
}
