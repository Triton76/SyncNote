package logic

import (
	"context"
	"errors"

	"SyncNote/rebuild/authapi/internal/svc"
	"SyncNote/rebuild/authapi/internal/types"
	"SyncNote/rebuild/common/model"

	"SyncNote/rebuild/pkg/crypto"
	"SyncNote/rebuild/pkg/util"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	//这里同时插入了两张表user和userinfo
	if req.Captcha != "" {
		return nil, errors.New("doesn't support captcha now")
	}
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("missing params")
	}

	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		logx.Errorf("register failed: hash password failed: %v", err)
		return nil, errors.New("system busy")
	}

	user := &model.User{
		UserId:       uuid.NewString(),
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	userInfo := &model.UserInfo{
		UserId:   user.UserId,
		Username: util.GenerateReadableUsername(),
	}

	//让ai帮忙改成事务插入了。 原先是两表单独插入，不保证数据一致性。
	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		txUserModel := model.NewUserModel(sqlx.NewSqlConnFromSession(session))
		txUserInfoModel := model.NewUserInfoModel(sqlx.NewSqlConnFromSession(session))

		if _, txErr := txUserModel.Insert(ctx, user); txErr != nil {
			return txErr
		}

		if _, txErr := txUserInfoModel.Insert(ctx, userInfo); txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, errors.New("email has been registered")
		}

		logx.Errorf("register failed: transaction failed: %v", err)
		return nil, errors.New("system busy")
	}

	return &types.RegisterResp{UserId: user.UserId}, nil
}
