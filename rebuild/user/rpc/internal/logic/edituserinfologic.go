package logic

import (
	"context"

	"SyncNote/rebuild/user/rpc/internal/svc"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserInfoLogic {
	return &EditUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditUserInfoLogic) EditUserInfo(in *userrpc.EditUserInfoReq) (*userrpc.Empty, error) {
	// todo: add your logic here and delete this line

	return &userrpc.Empty{}, nil
}
