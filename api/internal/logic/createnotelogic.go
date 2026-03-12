// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"SyncNote/api/internal/svc"
	"SyncNote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoteLogic {
	return &CreateNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNoteLogic) CreateNote(req *types.CreateNoteReq) (resp *types.NoteResp, err error) {
	// todo: add your logic here and delete this line

	return
}
