package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"SyncNote/syncnote/rpc/internal/model"
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取笔记的所有权限列表
func (l *ListPermissionsLogic) ListPermissions(in *syncnoterpc.ListPermissionsReq) (*syncnoterpc.ListPermissionsResp, error) {
	if in == nil {
		return nil, errors.New("request is nil")
	}
	if in.NoteId == "" || in.OperatorId == "" {
		return nil, errors.New("noteId and operatorId are required")
	}

	operator, err := l.svcCtx.NotePermissionsModel.FindOneByNoteIdUserId(l.ctx, in.NoteId, sql.NullString{String: in.OperatorId, Valid: true})
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("permission denied")
		}
		return nil, err
	}
	operatorRole := strings.ToLower(operator.Role)
	if operatorRole != "owner" && operatorRole != "admin" {
		return nil, errors.New("permission denied")
	}

	items, err := l.svcCtx.NotePermissionsModel.ListByNoteId(l.ctx, in.NoteId)
	if err != nil {
		return nil, err
	}

	permissions := make([]*syncnoterpc.PermissionInfo, 0, len(items))
	for _, item := range items {
		permissions = append(permissions, &syncnoterpc.PermissionInfo{
			PermissionId: item.PermissionId,
			NoteId:       item.NoteId,
			UserId:       item.UserId.String,
			TeamId:       item.TeamId.String,
			GrantedBy:    item.GrantedBy,
			Role:         dbToRole(item.Role),
			Status:       dbToStatus(item.Status),
			GrantedAt:    item.GrantedAt,
			RevokedAt:    item.RevokedAt.Int64,
		})
	}

	return &syncnoterpc.ListPermissionsResp{Permissions: permissions}, nil
}
