package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TeamModel = (*customTeamModel)(nil)

type (
	// TeamModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTeamModel.
	TeamModel interface {
		teamModel
		withSession(session sqlx.Session) TeamModel
	}

	customTeamModel struct {
		*defaultTeamModel
	}
)

// NewTeamModel returns a model for the database table.
func NewTeamModel(conn sqlx.SqlConn) TeamModel {
	return &customTeamModel{
		defaultTeamModel: newTeamModel(conn),
	}
}

func (m *customTeamModel) withSession(session sqlx.Session) TeamModel {
	return NewTeamModel(sqlx.NewSqlConnFromSession(session))
}
