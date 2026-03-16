// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"

	"SyncNote/auth/api/internal/model"
	"SyncNote/auth/api/internal/svc"
	"SyncNote/auth/api/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	user := &model.Users{
		Id:           newID.String(),
		PasswordHash: req.Password,
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
		Username:     sql.NullString{String: req.Username, Valid: req.Username != ""},
		Status:       1,
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}
	return &types.RegisterResp{
		Code:   model.CodeSuccess,
		UserId: user.Id,
		Token:  "233",
		Expire: 114514,
	}, nil
}
