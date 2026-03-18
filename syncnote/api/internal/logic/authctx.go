package logic

import (
	"context"
	"errors"

	"SyncNote/syncnote/api/internal/middleware"
)

func currentUserIDFromCtx(ctx context.Context) (string, error) {
	v := ctx.Value(middleware.CtxUserIDKey)
	userID, ok := v.(string)
	if !ok || userID == "" {
		return "", errors.New("unauthorized: user id missing in context")
	}
	return userID, nil
}
