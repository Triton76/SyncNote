package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TeamMembersModel = (*customTeamMembersModel)(nil)

type (
	// TeamMembersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTeamMembersModel.
	TeamMembersModel interface {
		teamMembersModel
		withSession(session sqlx.Session) TeamMembersModel
	}

	customTeamMembersModel struct {
		*defaultTeamMembersModel
	}
)

// NewTeamMembersModel returns a model for the database table.
func NewTeamMembersModel(conn sqlx.SqlConn) TeamMembersModel {
	return &customTeamMembersModel{
		defaultTeamMembersModel: newTeamMembersModel(conn),
	}
}

func (m *customTeamMembersModel) withSession(session sqlx.Session) TeamMembersModel {
	return NewTeamMembersModel(sqlx.NewSqlConnFromSession(session))
}
