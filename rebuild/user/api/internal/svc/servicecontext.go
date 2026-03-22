package svc

import (
	"SyncNote/rebuild/user/api/internal/config"
	"SyncNote/rebuild/user/rpc/userservice"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	userRpcClient := zrpc.MustNewClient(c.UserRpc)

	return &ServiceContext{
		Config:  c,
		UserRpc: userservice.NewUserService(userRpcClient),
	}
}
