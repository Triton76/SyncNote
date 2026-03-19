package logic

import (
	"context"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionsGetUserNotesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionsGetUserNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionsGetUserNotesLogic {
	return &OptionsGetUserNotesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionsGetUserNotesLogic) OptionsGetUserNotes(req *types.UserNotesReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
