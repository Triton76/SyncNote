package logic

import (
	"SyncNote/syncnote/rpc/internal/model"
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
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

func (l *CreateNoteLogic) CreateNote(in *syncnoterpc.CreateNoteReq) (*syncnoterpc.NoteResp, error) {
	if in.GetUserId() == "" {
		return nil, errors.New("userId is required")
	}
	if in.GetTitle() == "" {
		return nil, errors.New("title is required")
	}
	now := time.Now().UnixMilli()
	noteID := strings.ReplaceAll(uuid.NewString(), "-", "")
	note := &model.Notes{
		NoteId:       noteID,
		UserId:       in.UserId,
		Title:        in.Title,
		Content:      in.Content,
		Version:      1,
		LastModified: now,
	}
	res, err := l.svcCtx.NotesModel.Insert(l.ctx, note)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("save note failed")
	}

	ownerPermission := &model.NotePermissions{
		PermissionId: strings.ReplaceAll(uuid.NewString(), "-", ""),
		NoteId:       note.NoteId,
		UserId:       sql.NullString{String: in.UserId, Valid: true},
		TeamId:       sql.NullString{},
		GrantedBy:    in.UserId,
		Role:         "owner",
		Status:       "active",
		GrantedAt:    now,
		RevokedAt:    sql.NullInt64{},
	}
	if _, err = l.svcCtx.NotePermissionsModel.Insert(l.ctx, ownerPermission); err != nil {
		return nil, err
	}
	version := int64(note.Version)
	if err = appendCollaborationEvent(l.ctx, l.svcCtx, note.NoteId, "note_created", in.UserId, "", &version, nil, false); err != nil {
		return nil, err
	}

	return &syncnoterpc.NoteResp{
		NoteId:       note.NoteId,
		UserId:       note.UserId,
		Title:        note.Title,
		Content:      note.Content,
		Version:      int64(note.Version),
		LastModified: note.LastModified,
	}, nil
}
