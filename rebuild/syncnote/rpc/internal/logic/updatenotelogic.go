package logic

import (
	"context"
	"database/sql"
	"errors"

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
	// 乐观锁实现：当传递的Version低于最新Version，拒绝，当传递的Version等于当前数据库中的Version，同意并version+1
	// 为什么不用悲观锁？ 悲观锁事务开销大，要求将查询笔记、鉴权更新等步骤包裹在一个事务中，其中某个步骤耗时数据库就容易被占用。这时候还有一个性能问题，万一一个多人协作的笔记一人阻塞，导致后面所有人都无法进行编辑。
	//
	if in.GetNoteId() == "" {
		return nil, status.Error(codes.InvalidArgument, "note_id is required")
	}
	if in.GetVersion() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "version is required")
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

	err = l.svcCtx.NoteModel.UpdateWithVersion(l.ctx, note, int64(in.GetVersion()))
	if err != nil {
		if errors.Is(err, model.ErrOptimisticLockFailed) {
			return nil, status.Error(codes.Aborted, "version conflict")
		}
		return nil, err
	}

	latest, err := l.svcCtx.NoteModel.FindOne(l.ctx, in.GetNoteId())
	if err != nil {
		return nil, err
	}

	return &syncnoterpc.UpdateNoteResponse{Note: &syncnoterpc.Note{
		NoteId:  latest.NoteId,
		OwnerId: latest.OwnerId,
		Title:   latest.Title,
		Content: latest.Content.String,
		Version: int32(latest.Version),
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
