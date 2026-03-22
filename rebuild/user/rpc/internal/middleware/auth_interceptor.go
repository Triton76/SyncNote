package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthInterceptor gRPC 认证拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 从 metadata 提取 user_id
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		userIDs := md.Get("user_id")
		if len(userIDs) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing user_id")
		}

		// 注入到 context
		ctx = context.WithValue(ctx, UserIDKey, userIDs[0])

		return handler(ctx, req)
	}
}

// GetUserFromContext 获取用户 ID（Logic 层使用）
func GetUserFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}
	return userID, nil
}
