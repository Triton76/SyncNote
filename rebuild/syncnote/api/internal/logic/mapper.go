package logic

import (
	"time"

	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"
)

func toAPINote(note *syncnoterpc.Note) types.Note {
	if note == nil {
		return types.Note{}
	}

	createdAt := ""
	if note.GetCreatedAt() != nil {
		createdAt = note.GetCreatedAt().AsTime().Format(time.RFC3339)
	}

	updatedAt := ""
	if note.GetUpdatedAt() != nil {
		updatedAt = note.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}

	return types.Note{
		NoteId:    note.GetNoteId(),
		OwnerId:   note.GetOwnerId(),
		Title:     note.GetTitle(),
		Content:   note.GetContent(),
		Version:   note.GetVersion(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func toAPITeam(team *syncnoterpc.Team) types.Team {
	if team == nil {
		return types.Team{}
	}

	createdAt := ""
	if team.GetCreatedAt() != nil {
		createdAt = team.GetCreatedAt().AsTime().Format(time.RFC3339)
	}

	updatedAt := ""
	if team.GetUpdatedAt() != nil {
		updatedAt = team.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}

	return types.Team{
		TeamId:      team.GetTeamId(),
		Name:        team.GetName(),
		Description: team.GetDescription(),
		OwnerId:     team.GetOwnerId(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func toAPIUserPermission(permission *syncnoterpc.UserPermission) types.UserPermission {
	if permission == nil {
		return types.UserPermission{}
	}

	createdAt := ""
	if permission.GetCreatedAt() != nil {
		createdAt = permission.GetCreatedAt().AsTime().Format(time.RFC3339)
	}

	return types.UserPermission{
		PermissionId: permission.GetPermissionId(),
		NoteId:       permission.GetNoteId(),
		UserId:       permission.GetUserId(),
		Level:        int32(permission.GetLevel()),
		GrantedBy:    permission.GetGrantedBy(),
		CreatedAt:    createdAt,
	}
}