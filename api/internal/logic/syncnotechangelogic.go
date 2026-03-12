// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"SyncNote/api/internal/svc"
	"SyncNote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncNoteChangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncNoteChangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncNoteChangeLogic {
	return &SyncNoteChangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncNoteChangeLogic) SyncNoteChange(req *types.NoteChangeReq) (resp *types.SyncResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
