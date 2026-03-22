package logic

import (
	"context"
	"errors"

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
	userId := l.ctx.Value("userId")
	if userId == nil {
		return nil, errors.New("user not authenticated")
	}
	userIdStr, ok := userId.(string)
	if !ok || userIdStr == "" {
		return nil, errors.New("invalid user id")
	}

	getUserInfoReq := &userrpc.GetUserInfoReq{
		UserId: userIdStr,
	}
	inforesp, err := l.svcCtx.UserRpc.GetUserInfoById(l.ctx, getUserInfoReq)
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
