package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新笔记
func NewUpdateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoteLogic {
	return &UpdateNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNoteLogic) UpdateNote(req *types.UpdateNoteRequest) (resp *types.UpdateNoteResponse, err error) {
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.UpdateNote(rpcCtx, &syncnoterpc.UpdateNoteRequest{
		NoteId:  req.NoteId,
		Title:   req.Title,
		Content: req.Content,
		Version: req.Version,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateNoteResponse{Note: toAPINote(rpcResp.GetNote())}, nil
}
