package logic

import (
	"context"
	"errors"

	"SyncNote/rebuild/user/rpc/internal/svc"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserInfoByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByEmailLogic {
	return &GetUserInfoByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoByEmailLogic) GetUserInfoByEmail(in *userrpc.GetUserInfoByEmailReq) (*userrpc.GetUserInfoResp, error) {
	if in.Email == "" {
		return nil, errors.New("email required")
	}

	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	userInfo, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, user.UserId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &userrpc.GetUserInfoResp{
		UserInfo: &userrpc.UserInfo{
			Username:  userInfo.Username,
			UserId:    userInfo.UserId,
			Synopsis:  userInfo.Synopsis.String,
			AvatarUrl: userInfo.AvatarUrl.String,
			Email:     user.Email,
		},
	}, nil
}
