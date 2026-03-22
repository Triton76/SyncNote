package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoteLogic {
	return &DeleteNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNoteLogic) DeleteNote(in *syncnoterpc.DeleteNoteRequest) (*syncnoterpc.DeleteNoteResponse, error) {
	// 目前还是硬删除逻辑,要做鉴权（先查user permission，若user permission不符合再查团队权限，若都不符合拒绝）。
	
	return &syncnoterpc.DeleteNoteResponse{}, nil
}
