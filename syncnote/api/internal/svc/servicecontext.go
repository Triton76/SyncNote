// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"SyncNote/syncnote/api/internal/config"
	"SyncNote/syncnote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	SyncNoteRpc syncnoterpcclient.Syncnoterpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		SyncNoteRpc: syncnoterpcclient.NewSyncnoterpc(zrpc.MustNewClient(c.SyncNoteRpc)),
	}
}
