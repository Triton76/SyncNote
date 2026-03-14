// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"SyncNote/api/internal/svc"
	"SyncNote/api/internal/types"
	"SyncNote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteLogic {
	return &GetNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoteLogic) GetNote(req *types.NoteReq) (resp *types.NoteResp, err error) {
	// todo: add your logic here and delete this line
	rpcResp, err := l.svcCtx.SyncNoteRpc.GetNote(l.ctx, &syncnoterpcclient.NoteReq{
		NoteId: req.NoteId,
	})
	if err != nil {
		return nil, err
	}
	return &types.NoteResp{
		NoteId:       rpcResp.NoteId,
		UserId:       rpcResp.UserId,
		Title:        rpcResp.Title,
		Content:      rpcResp.Content,
		Version:      rpcResp.Version,
		LastModified: rpcResp.LastModified,
	}, nil
}
