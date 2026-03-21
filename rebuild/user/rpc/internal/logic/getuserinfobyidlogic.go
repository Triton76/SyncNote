package logic

import (
	"context"

	"SyncNote/rebuild/user/rpc/internal/svc"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByIdLogic {
	return &GetUserInfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoByIdLogic) GetUserInfoById(in *userrpc.GetUserInfoReq) (*userrpc.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &userrpc.GetUserInfoResp{}, nil
}
