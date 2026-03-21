package logic

import (
	"context"
	"errors"

	"SyncNote/rebuild/authapi/internal/svc"
	"SyncNote/rebuild/authapi/internal/types"
	"SyncNote/rebuild/pkg/auth"
	"SyncNote/rebuild/pkg/crypto"

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
	if req.Captcha != "" {
		return nil, errors.New("doesn't support captcha now")
	}
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		return nil, errors.New("login failed")
	}
	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("password wrong")
	}
	token, err := auth.GenerateToken(user.UserId)
	if err != nil {
		return nil, errors.New("token generate failed")
	}
	return &types.LoginResp{
		Token:    token,
		ExpireIn: int64(auth.TokenExpireDur.Seconds()),
	}, nil
}
