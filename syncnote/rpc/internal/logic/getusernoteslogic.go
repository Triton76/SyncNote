package logic

import (
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNotesLogic {
	return &GetUserNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserNotesLogic) GetUserNotes(in *syncnoterpc.UserNotesReq) (*syncnoterpc.UserNotesResp, error) {
	// 获取对应用户所有的Notes
	if in.GetUserId() == "" {
		return nil, errors.New("userId is required")
	}

	notesList, err := l.svcCtx.NotesModel.DumpUserNotes(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	//接下来进行Model层转化为proto层
	//提前分配内存节省性能。
	respNotes := make([]*syncnoterpc.NoteSummary, 0, len(notesList))

	for _, note:= range notesList {
		
		//在Model层用[]*Notes就记得需要判空，否则不需要下面这段判空代码。
		if note == nil {
			continue
		}

		summary := &syncnoterpc.NoteSummary {
			NoteId: note.NoteId,
			Title: note.Title,
			Version: int64(note.Version),
			LastModified: note.LastModified,
		}

		respNotes = append(respNotes, summary)
	}

	return &syncnoterpc.UserNotesResp {
		Notes: respNotes,
	}, nil
}
