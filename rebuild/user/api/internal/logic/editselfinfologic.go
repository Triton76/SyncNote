package logic

import (
	"context"

	"SyncNote/rebuild/user/api/internal/svc"
	"SyncNote/rebuild/user/api/internal/types"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditSelfInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditSelfInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditSelfInfoLogic {
	return &EditSelfInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditSelfInfoLogic) EditSelfInfo(req *types.EditReq) (resp *types.EmptyResp, err error) {
	editReq := &userrpc.EditUserInfoReq{Username: req.Username,
		Synopsis:  req.Synopsis,
		AvatarUrl: req.AvatarUrl,
	}
	_, err = l.svcCtx.UserRpc.EditUserInfo(l.ctx, editReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
