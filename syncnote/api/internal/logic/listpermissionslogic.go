package logic

import (
	"context"
	"errors"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"SyncNote/syncnote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPermissionsLogic) ListPermissions(req *types.ListPermissionsReq) (resp *types.ListPermissionsResp, err error) {
	operatorID, err := currentUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if req.NoteId == "" {
		return nil, errors.New("noteId is required")
	}

	rpcResp, err := l.svcCtx.SyncNoteRpc.ListPermissions(l.ctx, &syncnoterpcclient.ListPermissionsReq{
		NoteId:     req.NoteId,
		OperatorId: operatorID,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return &types.ListPermissionsResp{}, nil
	}

	items := make([]types.PermissionInfo, 0, len(rpcResp.Permissions))
	for _, p := range rpcResp.Permissions {
		items = append(items, types.PermissionInfo{
			PermissionId: p.PermissionId,
			NoteId:       p.NoteId,
			UserId:       p.UserId,
			TeamId:       p.TeamId,
			GrantedBy:    p.GrantedBy,
			Role:         roleToString(p.Role),
			Status:       permissionStatusToString(p.Status),
			GrantedAt:    p.GrantedAt,
			RevokedAt:    p.RevokedAt,
		})
	}
	return &types.ListPermissionsResp{Permissions: items}, nil
}
