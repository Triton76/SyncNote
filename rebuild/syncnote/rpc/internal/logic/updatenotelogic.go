package logic

import (
	"context"
	"database/sql"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoteLogic {
	return &UpdateNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateNoteLogic) UpdateNote(in *syncnoterpc.UpdateNoteRequest) (*syncnoterpc.UpdateNoteResponse, error) {
	// 更新笔记：owner 或有写权限以上的用户可更新。
	if in.GetNoteId() == "" {
		return nil, status.Error(codes.InvalidArgument, "note_id is required")
	}

	operatorID, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	note, err := l.svcCtx.NoteModel.FindOne(l.ctx, in.GetNoteId())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "note not found")
		}
		return nil, err
	}

	// 鉴权：owner 或写权限及以上
	if note.OwnerId != operatorID {
		ok, checkErr := l.canUpdateByPermission(in.GetNoteId(), operatorID)
		if checkErr != nil {
			return nil, checkErr
		}
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	// 更新笔记
	note.Title = in.GetTitle()
	note.Content = sql.NullString{String: in.GetContent(), Valid: in.GetContent() != ""}
	note.Version = int64(in.GetVersion())

	err = l.svcCtx.NoteModel.Update(l.ctx, note)
	if err != nil {
		return nil, err
	}

	return &syncnoterpc.UpdateNoteResponse{Note: &syncnoterpc.Note{
		NoteId:  note.NoteId,
		OwnerId: note.OwnerId,
		Title:   note.Title,
		Content: note.Content.String,
		Version: int32(note.Version),
	}}, nil
}

func (l *UpdateNoteLogic) canUpdateByPermission(noteId, operatorID string) (bool, error) {
	// 检查用户权限
	userPerm, err := l.svcCtx.NoteUserPermissionModel.FindOneByNoteIdUserId(l.ctx, noteId, operatorID)
	if err != nil && err != model.ErrNotFound {
		return false, err
	}
	if err == nil && model.HasPermissionLevel(userPerm.PermissionLevel, model.PermissionLevelWrite) {
		return true, nil
	}

	// 检查团队权限
	return l.svcCtx.NoteTeamPermissionModel.ExistsTeamPermissionLevel(l.ctx, noteId, operatorID, model.PermissionLevelWrite)
}
