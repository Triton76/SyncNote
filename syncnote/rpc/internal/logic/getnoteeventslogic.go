package logic

import (
	"context"

	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteEventsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoteEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteEventsLogic {
	return &GetNoteEventsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Collaboration History (对应 collaboration_events 表) ---
func (l *GetNoteEventsLogic) GetNoteEvents(in *syncnoterpc.GetNoteEventsReq) (*syncnoterpc.GetNoteEventsResp, error) {
	

	return &syncnoterpc.GetNoteEventsResp{}, nil
}
