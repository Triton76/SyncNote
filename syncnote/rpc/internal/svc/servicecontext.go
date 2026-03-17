package svc

import (
	"SyncNote/model"
	"SyncNote/rpc/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	NoteStore model.NoteStore
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		NoteStore: model.NewMockNoteStore(),
	}
}
