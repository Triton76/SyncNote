package logic

import (
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SaveNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoteLogic {
	return &SaveNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveNoteLogic) SaveNote(in *syncnoterpc.SaveNoteReq) (*syncnoterpc.SaveNoteResp, error) {
	//乐观锁：检查版本号，再更新，若版本号不匹配则拒绝更新，重新同步至最新的版本号后才接受更新。
	if in.GetNoteId() == "" {
		return nil, errors.New("note doesn't specified")
	}
	if in.GetUserId() == "" {
		return nil, errors.New("user_id is required")
	}
	note, err := l.svcCtx.NotesModel.FindOne(l.ctx, in.GetNoteId())
	if err == sqlx.ErrNotFound {
		return &syncnoterpc.SaveNoteResp{
			Success:       false,
			Code:          syncnoterpc.SaveCode_SAVE_CODE_NOT_FOUND,
			Message:       "Note not found.",
			LatestVersion: 0,
		}, nil
	}
	if err != nil {
		logx.Errorf("FindOne error: %v", err)
		return nil, err
	}
	//权限校验
	if note.UserId != in.UserId {
		return &syncnoterpc.SaveNoteResp{
			Success: false,
			Code:    syncnoterpc.SaveCode_SAVE_CODE_UNSPECIFIED, // 确保 proto 里有这个枚举
			Message: "Permission denied: cannot modify other user's note.",
		}, nil
	}
	if in.ExpectedVersion != int64(note.Version) {
		expected := in.ExpectedVersion
		latest := int64(note.Version)
		if logErr := appendCollaborationEvent(l.ctx, l.svcCtx, note.NoteId, "conflict_detected", in.UserId, "", &latest, &expected, true); logErr != nil {
			return nil, logErr
		}
		//这里是版本冲突，返回最新内容让客户端决定处理。
		return &syncnoterpc.SaveNoteResp{
			Success: false,
			Code:    syncnoterpc.SaveCode_SAVE_CODE_VERSION_CONFLICT,
			Message: "Save failed: version conflict.",
			Note: &syncnoterpc.NoteResp{
				NoteId:       note.NoteId,
				UserId:       note.UserId,
				Title:        note.Title,
				Content:      note.Content,
				Version:      int64(note.Version),
				LastModified: note.LastModified,
			},
			LatestVersion: int64(note.Version),
		}, nil
	}
	note.Version++
	if in.Title != "" {
		note.Title = in.Title
	}
	note.Content = in.Content
	note.LastModified = time.Now().UnixMilli()

	err = l.svcCtx.NotesModel.Update(l.ctx, note)
	if err != nil {
		logx.Errorf("Update note failed: %v", err)
		return nil, err
	}
	newVersion := int64(note.Version)
	if err = appendCollaborationEvent(l.ctx, l.svcCtx, note.NoteId, "note_updated", in.UserId, "", &newVersion, nil, false); err != nil {
		return nil, err
	}
	return &syncnoterpc.SaveNoteResp{
		Success: true,
		Code:    syncnoterpc.SaveCode_SAVE_CODE_OK,
		Message: "Save successed.",
		Note: &syncnoterpc.NoteResp{
			NoteId:       note.NoteId,
			UserId:       note.UserId,
			Title:        note.Title,
			Content:      note.Content,
			Version:      int64(note.Version),
			LastModified: note.LastModified,
		},
	}, nil

}
