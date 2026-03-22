package logic

import (
	"context"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/pkg/middleware"
	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

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

func (l *CreateNoteLogic) CreateNote(in *syncnoterpc.CreateNoteRequest) (*syncnoterpc.CreateNoteResponse, error) {
	// 创建，插入数据库，设置owner为当前用户。
	userId, err := middleware.GetUserFromContext(l.ctx)
		if err != nil {
		return nil, err
	}
	res, err := l.svcCtx.NoteModel.Insert(l.ctx, &model.Note{
		NoteId: uuid.NewString(),
		OwnerId: 
	})
	return &syncnoterpc.CreateNoteResponse{}, nil
}
