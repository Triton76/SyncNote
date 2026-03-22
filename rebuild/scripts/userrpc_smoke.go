package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func getenvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	rpcAddr := getenvOrDefault("RPC_ADDR", "127.0.0.1:8080")
	userID := os.Getenv("USER_ID")
	userEmail := os.Getenv("USER_EMAIL")
	newUsername := getenvOrDefault("NEW_USERNAME", fmt.Sprintf("rpc_tester_%d", time.Now().Unix()))

	if userID == "" {
		fmt.Fprintln(os.Stderr, "USER_ID is required")
		os.Exit(1)
	}
	if userEmail == "" {
		fmt.Fprintln(os.Stderr, "USER_EMAIL is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial %s failed: %v\n", rpcAddr, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := userrpc.NewUserServiceClient(conn)
	authCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	infoResp, err := client.GetUserInfoById(authCtx, &userrpc.GetUserInfoReq{UserId: userID})
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUserInfoById failed: %v\n", err)
		os.Exit(1)
	}
	if infoResp.GetUserInfo() == nil || infoResp.GetUserInfo().GetUserId() != userID {
		fmt.Fprintf(os.Stderr, "GetUserInfoById mismatch: got userId=%q\n", infoResp.GetUserInfo().GetUserId())
		os.Exit(1)
	}

	_, err = client.EditUserInfo(authCtx, &userrpc.EditUserInfoReq{
		Username:  newUsername,
		Synopsis:  "updated by userrpc smoke",
		AvatarUrl: "https://example.com/rpc-avatar.png",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "EditUserInfo failed: %v\n", err)
		os.Exit(1)
	}

	infoAfterEdit, err := client.GetUserInfoById(authCtx, &userrpc.GetUserInfoReq{UserId: userID})
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUserInfoById after edit failed: %v\n", err)
		os.Exit(1)
	}
	if infoAfterEdit.GetUserInfo() == nil || infoAfterEdit.GetUserInfo().GetUsername() != newUsername {
		fmt.Fprintf(os.Stderr, "Edit verification failed: expected username=%q, got=%q\n", newUsername, infoAfterEdit.GetUserInfo().GetUsername())
		os.Exit(1)
	}

	byEmailResp, err := client.GetUserInfoByEmail(authCtx, &userrpc.GetUserInfoByEmailReq{Email: userEmail})
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUserInfoByEmail failed: %v\n", err)
		os.Exit(1)
	}
	if byEmailResp.GetUserInfo() == nil || byEmailResp.GetUserInfo().GetUserId() != userID {
		fmt.Fprintf(os.Stderr, "GetUserInfoByEmail mismatch: expected userId=%q got=%q\n", userID, byEmailResp.GetUserInfo().GetUserId())
		os.Exit(1)
	}

	fmt.Println("userrpc smoke test passed")
	fmt.Printf("userId=%s\n", userID)
	fmt.Printf("email=%s\n", userEmail)
}
