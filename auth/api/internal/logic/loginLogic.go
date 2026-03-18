// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"
	"time"

	"SyncNote/auth/api/internal/model"
	"SyncNote/auth/api/internal/svc"
	"SyncNote/auth/api/internal/types"
	"SyncNote/pkg/auth"
	"SyncNote/pkg/crypto"

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
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, sql.NullString{String: req.LoginId, Valid: req.LoginId != ""})
	if err != nil {
		return nil, err
	}
	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		return &types.LoginResp{
			Code:   model.CodePasswordWrong,
			UserId: user.Id,
			Token:  "",
			Expire: 1,
		}, nil
	}
	token, err := auth.GenerateToken(req.LoginId)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		Code:   model.CodeSuccess,
		UserId: user.Id,
		Token:  token,
		Expire: time.Now().Add(auth.TokenExpireDur).Unix(),
	}, nil
}
