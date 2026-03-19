package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func mustEnv(name, fallback string) string {
	v := os.Getenv(name)
	if v == "" {
		return fallback
	}
	return v
}

func main() {
	rpcAddr := mustEnv("RPC_ADDR", "127.0.0.1:8080")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, rpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("dial rpc failed: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := syncnoterpc.NewSyncnoterpcClient(conn)
	userID := fmt.Sprintf("u%d%d", time.Now().Unix()%1000000, rand.Intn(1000))
	targetUserID := fmt.Sprintf("v%d%d", time.Now().Unix()%1000000, rand.Intn(1000))

	fmt.Println("[1/9] CreateNote...")
	createResp, err := client.CreateNote(ctx, &syncnoterpc.CreateNoteReq{
		UserId:  userID,
		Title:   "rpc smoke title",
		Content: "rpc smoke content",
	})
	if err != nil {
		log.Fatalf("CreateNote failed: %v", err)
	}
	if createResp.GetNoteId() == "" {
		log.Fatalf("CreateNote returned empty noteId")
	}
	noteID := createResp.GetNoteId()
	version := createResp.GetVersion()
	fmt.Printf("Created note_id=%s version=%d\n", noteID, version)

	fmt.Println("[2/9] GetNote...")
	getResp, err := client.GetNote(ctx, &syncnoterpc.NoteReq{NoteId: noteID})
	if err != nil {
		log.Fatalf("GetNote failed: %v", err)
	}
	if getResp.GetNoteId() != noteID {
		log.Fatalf("GetNote mismatch: expected %s, got %s", noteID, getResp.GetNoteId())
	}

	fmt.Println("[3/9] SaveNote success path...")
	saveResp, err := client.SaveNote(ctx, &syncnoterpc.SaveNoteReq{
		NoteId:          noteID,
		UserId:          userID,
		Content:         "rpc updated content",
		ExpectedVersion: version,
	})
	if err != nil {
		log.Fatalf("SaveNote failed: %v", err)
	}
	if !saveResp.GetSuccess() {
		log.Fatalf("SaveNote expected success=true, got false, code=%s msg=%s", saveResp.GetCode().String(), saveResp.GetMessage())
	}

	fmt.Println("[4/9] SaveNote conflict path...")
	conflictResp, err := client.SaveNote(ctx, &syncnoterpc.SaveNoteReq{
		NoteId:          noteID,
		UserId:          userID,
		Content:         "rpc conflict content",
		ExpectedVersion: version,
	})
	if err != nil {
		log.Fatalf("SaveNote conflict call failed: %v", err)
	}
	if conflictResp.GetCode() != syncnoterpc.SaveCode_SAVE_CODE_VERSION_CONFLICT {
		log.Fatalf("expected conflict code, got %s", conflictResp.GetCode().String())
	}

	fmt.Println("[5/9] GetUserNotes...")
	listResp, err := client.GetUserNotes(ctx, &syncnoterpc.UserNotesReq{UserId: userID})
	if err != nil {
		log.Fatalf("GetUserNotes failed: %v", err)
	}
	found := false
	for _, n := range listResp.GetNotes() {
		if n.GetNoteId() == noteID {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("GetUserNotes did not include noteId=%s", noteID)
	}

	fmt.Println("[6/9] GrantPermission...")
	grantResp, err := client.GrantPermission(ctx, &syncnoterpc.GrantPermissionReq{
		NoteId:       noteID,
		OperatorId:   userID,
		TargetUserId: targetUserID,
		Role:         syncnoterpc.Role_ROLE_EDITOR,
	})
	if err != nil {
		log.Fatalf("GrantPermission failed: %v", err)
	}
	if !grantResp.GetSuccess() {
		log.Fatalf("GrantPermission expected success=true, msg=%s", grantResp.GetMessage())
	}

	fmt.Println("[7/9] ListPermissions...")
	permResp, err := client.ListPermissions(ctx, &syncnoterpc.ListPermissionsReq{
		NoteId:     noteID,
		OperatorId: userID,
	})
	if err != nil {
		log.Fatalf("ListPermissions failed: %v", err)
	}
	permFound := false
	for _, p := range permResp.GetPermissions() {
		if p.GetUserId() == targetUserID && p.GetStatus() == syncnoterpc.PermissionStatus_PERMISSION_STATUS_ACTIVE {
			permFound = true
			break
		}
	}
	if !permFound {
		log.Fatalf("ListPermissions did not include active target permission")
	}

	fmt.Println("[8/9] RevokePermission...")
	revokeResp, err := client.RevokePermission(ctx, &syncnoterpc.RevokePermissionReq{
		NoteId:       noteID,
		OperatorId:   userID,
		TargetUserId: targetUserID,
	})
	if err != nil {
		log.Fatalf("RevokePermission failed: %v", err)
	}
	if !revokeResp.GetSuccess() {
		log.Fatalf("RevokePermission expected success=true, msg=%s", revokeResp.GetMessage())
	}

	fmt.Println("[9/9] GetNoteEvents...")
	eventsResp, err := client.GetNoteEvents(ctx, &syncnoterpc.GetNoteEventsReq{
		NoteId:   noteID,
		StartSeq: 0,
		Limit:    20,
	})
	if err != nil {
		log.Fatalf("GetNoteEvents failed: %v", err)
	}
	if eventsResp.GetEvents() == nil {
		log.Fatalf("GetNoteEvents returned nil events")
	}

	fmt.Println("RPC smoke test passed.")
}
