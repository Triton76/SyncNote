package logic

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

func getUserIDFromContext(ctx context.Context) (string, error) {
	keys := []string{"userId", "user_id", "uid"}
	for _, key := range keys {
		value := ctx.Value(key)
		if userID, ok := value.(string); ok && userID != "" {
			return userID, nil
		}
	}

	return "", errors.New("user not authenticated")
}

func withRPCUserID(ctx context.Context, userID string) context.Context {
	if userID == "" {
		return ctx
	}

	return metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"userId", userID,
		"user_id", userID,
	))
}
