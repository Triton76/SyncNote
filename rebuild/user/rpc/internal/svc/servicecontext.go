package svc

import (
	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserModel     model.UserModel
	UserInfoModel model.UserInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	userModel := model.NewUserModel(conn)
	userInfoModel := model.NewUserInfoModel(conn)
	return &ServiceContext{
		Config:        c,
		UserModel:     userModel,
		UserInfoModel: userInfoModel,
	}
}
