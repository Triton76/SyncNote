package logic

import (
	"context"

	"SyncNote/rebuild/user/api/internal/svc"
	"SyncNote/rebuild/user/api/internal/types"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelfInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfInfoLogic {
	return &GetSelfInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelfInfoLogic) GetSelfInfo() (resp *types.GetUserInfoResp, err error) {
	userID, err := getUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	getUserInfoReq := &userrpc.GetUserInfoReq{
		UserId: userID,
	}
	rpcCtx := withRPCUserID(l.ctx, userID)
	inforesp, err := l.svcCtx.UserRpc.GetUserInfoById(rpcCtx, getUserInfoReq)
	if err != nil {
		return nil, err
	}
	return &types.GetUserInfoResp{
		UserInfo: types.UserInfo{
			Username:  inforesp.UserInfo.Username,
			UserId:    inforesp.UserInfo.UserId,
			Email:     inforesp.UserInfo.Email,
			Synopsis:  inforesp.UserInfo.Synopsis,
			AvatarUrl: inforesp.UserInfo.AvatarUrl,
		},
	}, nil
}
