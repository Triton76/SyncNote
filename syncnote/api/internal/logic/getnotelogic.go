// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"SyncNote/syncnote/rpc/syncnoterpcclient"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrForbiddenNoteAccess = errors.New("forbidden: note does not belong to current user")

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
	userID, err := currentUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.SyncNoteRpc.GetNote(l.ctx, &syncnoterpcclient.NoteReq{
		NoteId: req.NoteId,
	})
	if err != nil {
		return nil, err
	}

	if rpcResp.UserId != userID {
		return nil, ErrForbiddenNoteAccess
	} //这里新加了鉴权逻辑，不然一个受认证的用户A能直接访问用户B的note。

	return &types.NoteResp{
		NoteId:       rpcResp.NoteId,
		UserId:       rpcResp.UserId,
		Title:        rpcResp.Title,
		Content:      rpcResp.Content,
		Version:      rpcResp.Version,
		LastModified: rpcResp.LastModified,
	}, nil
}
