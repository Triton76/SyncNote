package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type noCacheUsersModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewUsersModelWithoutCache returns a model that talks to MySQL directly without Redis cache.
func NewUsersModelWithoutCache(conn sqlx.SqlConn) UsersModel {
	return &noCacheUsersModel{
		conn:  conn,
		table: "`users`",
	}
}

func (m *noCacheUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, usersRowsExpectAutoSet)
	return m.conn.ExecCtx(ctx, query, data.Id, data.PasswordHash, data.Email, data.Username, data.Status)
}

func (m *noCacheUsersModel) FindOne(ctx context.Context, id string) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &resp, nil
}

func (m *noCacheUsersModel) FindOneByEmail(ctx context.Context, email sql.NullString) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &resp, nil
}

func (m *noCacheUsersModel) Update(ctx context.Context, newData *Users) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, usersRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.PasswordHash, newData.Email, newData.Username, newData.Status, newData.Id)
	return err
}

func (m *noCacheUsersModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
