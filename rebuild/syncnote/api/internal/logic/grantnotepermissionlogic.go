package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"SyncNote/rebuild/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantNotePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 授予笔记权限
func NewGrantNotePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantNotePermissionLogic {
	return &GrantNotePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GrantNotePermissionLogic) GrantNotePermission(req *types.GrantNotePermissionRequest) (resp *types.GrantNotePermissionResponse, err error) {
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcCtx := withRPCUserID(l.ctx, userID)
	rpcResp, err := l.svcCtx.SyncnoteRpc.GrantNotePermission(rpcCtx, &syncnoterpc.GrantNotePermissionRequest{
		NoteId:       req.NoteId,
		TargetUserId: req.TargetUserId,
		Level:        syncnoterpc.PermissionLevel(req.Level),
	})
	if err != nil {
		return nil, err
	}

	return &types.GrantNotePermissionResponse{Permission: toAPIUserPermission(rpcResp.GetPermission())}, nil
}
