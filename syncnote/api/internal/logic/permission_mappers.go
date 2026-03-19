package logic

import (
	"strings"

	"SyncNote/syncnote/rpc/pb/syncnoterpc"
)

func parseRole(input string) (syncnoterpc.Role, bool) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "owner":
		return syncnoterpc.Role_ROLE_OWNER, true
	case "admin":
		return syncnoterpc.Role_ROLE_ADMIN, true
	case "editor":
		return syncnoterpc.Role_ROLE_EDITOR, true
	case "viewer":
		return syncnoterpc.Role_ROLE_VIEWER, true
	default:
		return syncnoterpc.Role_ROLE_UNSPECIFIED, false
	}
}

func roleToString(role syncnoterpc.Role) string {
	switch role {
	case syncnoterpc.Role_ROLE_OWNER:
		return "owner"
	case syncnoterpc.Role_ROLE_ADMIN:
		return "admin"
	case syncnoterpc.Role_ROLE_EDITOR:
		return "editor"
	case syncnoterpc.Role_ROLE_VIEWER:
		return "viewer"
	default:
		return "unspecified"
	}
}

func permissionStatusToString(status syncnoterpc.PermissionStatus) string {
	switch status {
	case syncnoterpc.PermissionStatus_PERMISSION_STATUS_ACTIVE:
		return "active"
	case syncnoterpc.PermissionStatus_PERMISSION_STATUS_REVOKED:
		return "revoked"
	case syncnoterpc.PermissionStatus_PERMISSION_STATUS_PENDING:
		return "pending"
	default:
		return "unspecified"
	}
}

func eventTypeToString(t syncnoterpc.EventType) string {
	switch t {
	case syncnoterpc.EventType_EVENT_TYPE_NOTE_CREATED:
		return "note_created"
	case syncnoterpc.EventType_EVENT_TYPE_NOTE_UPDATED:
		return "note_updated"
	case syncnoterpc.EventType_EVENT_TYPE_NOTE_DELETED:
		return "note_deleted"
	case syncnoterpc.EventType_EVENT_TYPE_PERMISSION_GRANTED:
		return "permission_granted"
	case syncnoterpc.EventType_EVENT_TYPE_PERMISSION_REVOKED:
		return "permission_revoked"
	case syncnoterpc.EventType_EVENT_TYPE_CONFLICT_DETECTED:
		return "conflict_detected"
	case syncnoterpc.EventType_EVENT_TYPE_VIEW_STARTED:
		return "view_started"
	case syncnoterpc.EventType_EVENT_TYPE_VIEW_ENDED:
		return "view_ended"
	default:
		return "unspecified"
	}
}
