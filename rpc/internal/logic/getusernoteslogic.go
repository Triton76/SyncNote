package logic

import (
	"context"

	"SyncNote/rpc/internal/svc"
	"SyncNote/rpc/pb/syncnoterpc"

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
	// todo: add your logic here and delete this line

	return &syncnoterpc.UserNotesResp{}, nil
}
