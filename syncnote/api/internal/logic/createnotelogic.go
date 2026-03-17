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

type CreateNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoteLogic {
	return &CreateNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNoteLogic) CreateNote(req *types.CreateNoteReq) (resp *types.NoteResp, err error) {
	// todo: add your logic here and delete this line
	rpcResp, err := l.svcCtx.SyncNoteRpc.CreateNote(l.ctx, &syncnoterpcclient.CreateNoteReq{
		UserId:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
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
