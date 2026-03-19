package logic

import (
	"context"
	"errors"

	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"SyncNote/syncnote/rpc/syncnoterpcclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGrantPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantPermissionLogic {
	return &GrantPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GrantPermissionLogic) GrantPermission(req *types.GrantPermissionReq) (resp *types.PermissionResp, err error) {
	operatorID, err := currentUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if req.NoteId == "" {
		return nil, errors.New("noteId is required")
	}
	role, ok := parseRole(req.Role)
	if !ok {
		return nil, errors.New("invalid role, expected one of owner/admin/editor/viewer")
	}

	rpcResp, err := l.svcCtx.SyncNoteRpc.GrantPermission(l.ctx, &syncnoterpcclient.GrantPermissionReq{
		NoteId:          req.NoteId,
		OperatorId:      operatorID,
		TargetUserId:    req.TargetUserId,
		TargetUserEmail: req.TargetUserEmail,
		TargetTeamId:    req.TargetTeamId,
		Role:            role,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return &types.PermissionResp{}, nil
	}

	out := &types.PermissionResp{Success: rpcResp.Success, Message: rpcResp.Message}
	if rpcResp.Permission != nil {
		out.Permission = &types.PermissionInfo{
			PermissionId: rpcResp.Permission.PermissionId,
			NoteId:       rpcResp.Permission.NoteId,
			UserId:       rpcResp.Permission.UserId,
			TeamId:       rpcResp.Permission.TeamId,
			GrantedBy:    rpcResp.Permission.GrantedBy,
			Role:         roleToString(rpcResp.Permission.Role),
			Status:       permissionStatusToString(rpcResp.Permission.Status),
			GrantedAt:    rpcResp.Permission.GrantedAt,
			RevokedAt:    rpcResp.Permission.RevokedAt,
		}
	}
	return out, nil
}
