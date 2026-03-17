// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"SyncNote/syncnote/rpc/syncnoterpcclient"
	"context"

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
	rpcResp, err := l.svcCtx.SyncNoteRpc.GetUserNotes(l.ctx, &syncnoterpcclient.UserNotesReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return &types.UserNotesResp{}, nil
	}
	notes := make([]types.NoteSummary, len(rpcResp.Notes))
	for i, note := range rpcResp.Notes {
		notes[i] = types.NoteSummary{
			NoteId:       note.NoteId,
			Title:        note.Title,
			Version:      note.Version,
			LastModified: note.LastModified,
		}
	}
	return &types.UserNotesResp{Notes: notes}, nil

}
