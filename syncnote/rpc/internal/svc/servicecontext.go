package svc

import (
	"SyncNote/syncnote/rpc/internal/config"
	"SyncNote/syncnote/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	Conn                 sqlx.SqlConn
	NotesModel           model.NotesModel
	NotePermissionsModel model.NotePermissionsModel
	CollabEventsModel    model.CollaborationEventsModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn := sqlx.NewMysql(c.DataSource)
	notesModel := model.NewNotesModel(conn, c.CacheRedis)
	notePermissionsModel := model.NewNotePermissionsModel(conn, c.CacheRedis)
	collabEventsModel := model.NewCollaborationEventsModel(conn, c.CacheRedis)
	return &ServiceContext{
		Config:               c,
		Conn:                 conn,
		NotesModel:           notesModel,
		NotePermissionsModel: notePermissionsModel,
		CollabEventsModel:    collabEventsModel,
	}
}
