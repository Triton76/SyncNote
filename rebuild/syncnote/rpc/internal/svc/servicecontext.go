package svc

import (
	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/syncnote/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	NoteModel               model.NoteModel
	NoteTeamPermissionModel model.NoteTeamPermissionModel
	NoteUserPermissionModel model.NoteUserPermissionModel
	TeamMembersModel        model.TeamMembersModel
	TeamModel               model.TeamModel
	UserModel               model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:                  c,
		NoteModel:               model.NewNoteModel(conn),
		NoteTeamPermissionModel: model.NewNoteTeamPermissionModel(conn),
		NoteUserPermissionModel: model.NewNoteUserPermissionModel(conn),
		TeamMembersModel:        model.NewTeamMembersModel(conn),
		TeamModel:               model.NewTeamModel(conn),
		UserModel:               model.NewUserModel(conn),
	}
}
