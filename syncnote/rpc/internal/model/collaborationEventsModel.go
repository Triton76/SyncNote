package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CollaborationEventsModel = (*customCollaborationEventsModel)(nil)

type (
	// CollaborationEventsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCollaborationEventsModel.
	CollaborationEventsModel interface {
		collaborationEventsModel
	}

	customCollaborationEventsModel struct {
		*defaultCollaborationEventsModel
	}
)

// NewCollaborationEventsModel returns a model for the database table.
func NewCollaborationEventsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CollaborationEventsModel {
	return &customCollaborationEventsModel{
		defaultCollaborationEventsModel: newCollaborationEventsModel(conn, c, opts...),
	}
}
