package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"SyncNote/rebuild/authapi/internal/svc"
	"SyncNote/rebuild/authapi/internal/types"
	"SyncNote/rebuild/common/model"

	"SyncNote/rebuild/pkg/crypto"

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
	if req.Captcha != "" {
		return nil, errors.New("doesn't support captcha now")
	}
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("missing params")
	}
	passwordHash, _ := crypto.HashPassword(req.Password)
	user := &model.User{
		UserId:       uuid.NewString(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		logx.Error(fmt.Sprintf("register failed: insert failed %v", err))
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, errors.New("email has been registered")
		}
		return nil, errors.New("system busy")
	}
	return &types.RegisterResp{UserId: user.UserId}, nil
}
