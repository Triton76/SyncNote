package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoteTeamPermissionModel = (*customNoteTeamPermissionModel)(nil)

type (
	// NoteTeamPermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoteTeamPermissionModel.
	NoteTeamPermissionModel interface {
		noteTeamPermissionModel
		ExistsTeamAdminForUser(ctx context.Context, noteId, userId string) (bool, error)
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

func (m *customNoteTeamPermissionModel) ExistsTeamAdminForUser(ctx context.Context, noteId, userId string) (bool, error) {
	query := fmt.Sprintf(`
		select count(1)
		from %s ntp
		join %s tm on tm.team_id = ntp.team_id
		join %s n on n.note_id = ntp.note_id
		join %s t on t.team_id = ntp.team_id
		where ntp.note_id = ?
		  and tm.user_id = ?
		  and ntp.permission_level = 'admin'
		  and n.deleted_at is null
		  and t.deleted_at is null
	`, m.table, "`team_members`", "`note`", "`team`")

	var cnt int64
	if err := m.conn.QueryRowCtx(ctx, &cnt, query, noteId, userId); err != nil {
		return false, err
	}

	return cnt > 0, nil
}
