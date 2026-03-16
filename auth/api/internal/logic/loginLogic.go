// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"

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
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, sql.NullString{String: req.LoginId, Valid: req.LoginId != ""})
	if err != nil {
		return nil, err
	}
	if user.PasswordHash != req.Password {
		return nil, err
	}
	return &types.LoginResp{
		Code:   model.CodeSuccess,
		UserId: user.Id,
		Token:  "233",
		Expire: 114514,
	}, nil
}
