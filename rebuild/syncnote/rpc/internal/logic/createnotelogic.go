package logic

import (
	"context"
	"database/sql"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoteLogic {
	return &CreateNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNoteLogic) CreateNote(in *syncnoterpc.CreateNoteRequest) (*syncnoterpc.CreateNoteResponse, error) {
	// 创建，插入数据库，设置owner为当前用户。
	if in.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	userId, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	content := sql.NullString{String: in.Content, Valid: in.Content != ""}
	note := &model.Note{
		NoteId:  uuid.NewString(),
		OwnerId: userId,
		Title:   in.Title,
		Content: content,
		Version: 1,
	}
	_, err = l.svcCtx.NoteModel.Insert(l.ctx, note)
	if err != nil {
		return nil, err
	}
	return &syncnoterpc.CreateNoteResponse{Note: &syncnoterpc.Note{
		NoteId:  note.NoteId,
		OwnerId: note.OwnerId,
		Title:   note.Title,
		Content: note.Content.String,
		Version: int32(note.Version),
	}}, nil
}
