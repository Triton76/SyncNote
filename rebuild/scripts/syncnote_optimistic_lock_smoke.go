package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func getenvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func failf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	rpcAddr := getenvOrDefault("RPC_ADDR", "127.0.0.1:8080")
	userID := os.Getenv("USER_ID")
	if userID == "" {
		failf("USER_ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		failf("dial %s failed: %v", rpcAddr, err)
	}
	defer conn.Close()

	client := syncnoterpc.NewSyncnoteServiceClient(conn)
	authCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	createResp, err := client.CreateNote(authCtx, &syncnoterpc.CreateNoteRequest{
		Title:   fmt.Sprintf("optimistic-lock-test-%d", time.Now().UnixNano()),
		Content: "initial",
	})
	if err != nil {
		failf("CreateNote failed: %v", err)
	}
	if createResp.GetNote() == nil {
		failf("CreateNote returned empty note")
	}

	noteID := createResp.GetNote().GetNoteId()
	baseVersion := createResp.GetNote().GetVersion()
	if noteID == "" || baseVersion <= 0 {
		failf("invalid note created: noteID=%q, version=%d", noteID, baseVersion)
	}

	_, err = client.UpdateNote(authCtx, &syncnoterpc.UpdateNoteRequest{
		NoteId:  noteID,
		Title:   "v2-first-update",
		Content: "first-update",
		Version: baseVersion,
	})
	if err != nil {
		failf("first UpdateNote failed unexpectedly: %v", err)
	}

	_, err = client.UpdateNote(authCtx, &syncnoterpc.UpdateNoteRequest{
		NoteId:  noteID,
		Title:   "v2-second-update-should-conflict",
		Content: "second-update",
		Version: baseVersion,
	})
	if err == nil {
		failf("second UpdateNote unexpectedly succeeded; optimistic lock did not trigger")
	}

	st, ok := status.FromError(err)
	if !ok {
		failf("second UpdateNote returned non-gRPC error: %v", err)
	}
	if st.Code() != codes.Aborted {
		failf("expected Aborted conflict, got %s (%v)", st.Code().String(), err)
	}

	// cleanup best effort
	_, _ = client.DeleteNote(authCtx, &syncnoterpc.DeleteNoteRequest{NoteId: noteID})

	fmt.Println("syncnote optimistic lock smoke test passed")
	fmt.Printf("noteId=%s\n", noteID)
	fmt.Printf("baseVersion=%d\n", baseVersion)
}
