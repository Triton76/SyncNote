package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeNotePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤销笔记权限
func NewRevokeNotePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeNotePermissionLogic {
	return &RevokeNotePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeNotePermissionLogic) RevokeNotePermission(req *types.RevokeNotePermissionRequest) (resp *types.RevokeNotePermissionResponse, err error) {
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.RevokeNotePermission(rpcCtx, &syncnoterpc.RevokeNotePermissionRequest{
		NoteId:       req.NoteId,
		TargetUserId: req.TargetUserId,
	})
	if err != nil {
		return nil, err
	}

	return &types.RevokeNotePermissionResponse{Success: rpcResp.GetSuccess()}, nil
}
