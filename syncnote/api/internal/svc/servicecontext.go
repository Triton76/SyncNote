// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"SyncNote/syncnote/api/internal/config"
	"SyncNote/syncnote/api/internal/middleware"
	"SyncNote/syncnote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	SyncNoteRpc syncnoterpcclient.Syncnoterpc
	Auth        rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	authMiddleware := middleware.NewAuthMiddleware()

	return &ServiceContext{
		Config:      c,
		SyncNoteRpc: syncnoterpcclient.NewSyncnoterpc(zrpc.MustNewClient(c.SyncNoteRpc)),
		Auth:        authMiddleware.Handle,
	}
}
