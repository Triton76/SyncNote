// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"SyncNote/api/internal/svc"
	"SyncNote/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNotesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNotesLogic {
	return &GetUserNotesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserNotesLogic) GetUserNotes(req *types.UserNotesReq) (resp *types.UserNotesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
