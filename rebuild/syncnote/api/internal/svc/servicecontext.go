package svc

import (
	"SyncNote/rebuild/syncnote/api/internal/config"
	"SyncNote/rebuild/syncnote/rpc/syncnoteservice"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	SyncnoteRpc syncnoteservice.SyncnoteService
}

func NewServiceContext(c config.Config) *ServiceContext {
	rpcClient := zrpc.MustNewClient(c.SyncnoteRpc)
	return &ServiceContext{
		Config:      c,
		SyncnoteRpc: syncnoteservice.NewSyncnoteService(rpcClient),
	}
}
