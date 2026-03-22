package logic

import (
	"context"
	"errors"

	"SyncNote/rebuild/user/api/internal/svc"
	"SyncNote/rebuild/user/api/internal/types"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.GetUserInfoResp, err error) {
	if req.UserId == "" {
		return nil, errors.New("userId required")
	}

	inforesp, err := l.svcCtx.UserRpc.GetUserInfoById(l.ctx, &userrpc.GetUserInfoReq{UserId: req.UserId})
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
