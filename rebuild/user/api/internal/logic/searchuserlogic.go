package logic

import (
	"context"
	"errors"

	"SyncNote/rebuild/user/api/internal/svc"
	"SyncNote/rebuild/user/api/internal/types"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchUserLogic) SearchUser(req *types.SearchReq) (resp *types.SearchResp, err error) {
	if req.Email == "" {
		return nil, errors.New("email required")
	}

	inforesp, err := l.svcCtx.UserRpc.GetUserInfoByEmail(l.ctx, &userrpc.GetUserInfoByEmailReq{Email: req.Email})
	if err != nil {
		return nil, err
	}

	return &types.SearchResp{
		InfoList: []types.UserInfo{
			{
				Username:  inforesp.UserInfo.Username,
				UserId:    inforesp.UserInfo.UserId,
				Email:     inforesp.UserInfo.Email,
				Synopsis:  inforesp.UserInfo.Synopsis,
				AvatarUrl: inforesp.UserInfo.AvatarUrl,
			},
		},
		Total: 1,
	}, nil
}
