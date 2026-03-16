// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"SyncNote/auth/api/internal/config"

	"SyncNote/auth/api/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config: c,

		UserModel: model.NewUsersModel(conn),
	}
}
