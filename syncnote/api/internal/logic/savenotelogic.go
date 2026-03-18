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

type SaveNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoteLogic {
	return &SaveNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveNoteLogic) SaveNote(req *types.SaveNoteReq) (resp *types.SaveNoteResp, err error) {
	rpcResp, err := l.svcCtx.SyncNoteRpc.SaveNote(l.ctx, &syncnoterpcclient.SaveNoteReq{
		NoteId:          req.NoteId,
		UserId:          req.UserId,
		Content:         req.Content,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return &types.SaveNoteResp{}, nil
	}
	var note *types.NoteResp
	if rpcResp.Note != nil {
		note = &types.NoteResp{
			NoteId:       rpcResp.Note.NoteId,
			UserId:       rpcResp.Note.UserId,
			Title:        rpcResp.Note.Title,
			Content:      rpcResp.Note.Content,
			Version:      rpcResp.Note.Version,
			LastModified: rpcResp.Note.LastModified,
		}
	}
	return &types.SaveNoteResp{
		Success:       rpcResp.Success,
		Code:          rpcResp.Code.String(),
		Message:       rpcResp.Message,
		Note:          note,
		LatestVersion: rpcResp.LatestVersion,
	}, nil
}
