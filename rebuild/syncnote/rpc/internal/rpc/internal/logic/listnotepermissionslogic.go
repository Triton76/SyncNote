package logic

import (
	"context"

	"SyncNote/rebuild/syncnote/rpc/internal/rpc/internal/svc"
	"SyncNote/rebuild/syncnote/rpc/internal/rpc/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotePermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNotePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotePermissionsLogic {
	return &ListNotePermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNotePermissionsLogic) ListNotePermissions(in *syncnoterpc.ListNotePermissionsRequest) (*syncnoterpc.ListNotePermissionsResponse, error) {
	// todo: add your logic here and delete this line

	return &syncnoterpc.ListNotePermissionsResponse{}, nil
}
