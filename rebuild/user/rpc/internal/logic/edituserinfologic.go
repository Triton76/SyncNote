package logic

import (
	"context"
	"database/sql"
	"fmt"

	"SyncNote/rebuild/common/model"
	"SyncNote/rebuild/user/rpc/internal/middleware"
	"SyncNote/rebuild/user/rpc/internal/svc"
	"SyncNote/rebuild/user/rpc/pb/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserInfoLogic {
	return &EditUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditUserInfoLogic) EditUserInfo(in *userrpc.EditUserInfoReq) (*userrpc.Empty, error) {
	userId, err := middleware.GetUserFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	synopsis := sql.NullString{String: in.Synopsis, Valid: in.Synopsis != ""}
	avatarUrl := sql.NullString{String: in.AvatarUrl, Valid: in.AvatarUrl != ""}
	userInfo := &model.UserInfo{
		UserId:    userId,
		Username:  in.Username,
		Synopsis:  synopsis,
		AvatarUrl: avatarUrl,
	}
	err = l.svcCtx.UserInfoModel.Update(l.ctx, userInfo)
	if err != nil {
		logx.Error(fmt.Sprintf("update failed: insert failed %v", err))
		return nil, err
	}
	return &userrpc.Empty{}, nil
}
