package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrOptimisticLockFailed = errors.New("optimistic lock failed")

// 权限级别常量
const (
	PermissionLevelRead  = "read"
	PermissionLevelWrite = "write"
	PermissionLevelAdmin = "admin"
)

// PermissionLevelValue 返回权限级别的数值（用于比较）
// read=1, write=2, admin=3
func PermissionLevelValue(level string) int {
	switch level {
	case PermissionLevelRead:
		return 1
	case PermissionLevelWrite:
		return 2
	case PermissionLevelAdmin:
		return 3
	default:
		return 0
	}
}

// HasPermissionLevel 判断实际权限是否满足要求的权限级别
// 例：HasPermissionLevel("write", "read") -> true（写权限满足读要求）
//     HasPermissionLevel("read", "write") -> false（读权限不满足写要求）
func HasPermissionLevel(actualLevel, requiredLevel string) bool {
	return PermissionLevelValue(actualLevel) >= PermissionLevelValue(requiredLevel)
}
