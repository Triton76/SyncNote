package svc

import (
	"SyncNote/rebuild/authapi/internal/config"
	"SyncNote/rebuild/common/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	userModel := model.NewUserModel(conn)
	return &ServiceContext{
		Config:    c,
		UserModel: userModel,
	}
}
