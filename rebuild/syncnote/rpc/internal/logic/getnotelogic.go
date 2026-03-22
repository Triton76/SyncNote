package logic

import (
	"context"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteLogic {
	return &GetNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoteLogic) GetNote(in *syncnoterpc.GetNoteRequest) (*syncnoterpc.GetNoteResponse, error) {
	// 鉴权，链条为： owner-> user_permission -> team_permission

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

	// 鉴权：1. owner 直接允许 2. user_permission 3. team_permission
	hasAccess := note.OwnerId == operatorID
	if !hasAccess {
		// 检查 user_permission
		userPerm, err := l.svcCtx.NoteUserPermissionModel.FindOneByNoteIdUserId(l.ctx, in.GetNoteId(), operatorID)
		if err != nil && err != model.ErrNotFound {
			return nil, err
		}
		if err == nil && userPerm.PermissionLevel != "" {
			hasAccess = true
		}
	}

	if !hasAccess {
		// 检查 team_permission（任何权限级别都可以读）
		teamAccess, err := l.svcCtx.NoteTeamPermissionModel.ExistsTeamPermissionForUser(l.ctx, in.GetNoteId(), operatorID)
		if err != nil {
			return nil, err
		}
		hasAccess = teamAccess
	}

	if !hasAccess {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	return &syncnoterpc.GetNoteResponse{Note: &syncnoterpc.Note{
		NoteId:  note.NoteId,
		OwnerId: note.OwnerId,
		Title:   note.Title,
		Content: note.Content.String,
		Version: int32(note.Version),
	}}, nil
}
