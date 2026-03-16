// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"SyncNote/auth/api/internal/model"
	"SyncNote/auth/api/internal/svc"
	"SyncNote/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 这里暂时使用Email作为登陆ID
	user, err := l.svcCtx.UserModel.FindOneByAccount(l.ctx, req.LoginId)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		Code:   model.CodeSuccess,
		UserId: strconv.FormatUint(user.Id, 10),
		Token:  "233test",
		Expire: 114514,
	}, nil
}
