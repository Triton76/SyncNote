package svc

import (
	"SyncNote/rebuild/authapi/internal/config"
	"SyncNote/rebuild/common/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	DB            sqlx.SqlConn
	UserModel     model.UserModel
	UserInfoModel model.UserInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	userModel := model.NewUserModel(conn)
	userInfoModel := model.NewUserInfoModel(conn)
	return &ServiceContext{
		Config:        c,
		DB:            conn,
		UserModel:     userModel,
		UserInfoModel: userInfoModel,
	}
}
