package logic

import (
	"context"
	"database/sql"
	"errors"

	"SyncNote/syncnote/rpc/internal/model"
	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GrantPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGrantPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantPermissionLogic {
	return &GrantPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Permission Management (对应 note_permissions 表) ---
func (l *GrantPermissionLogic) GrantPermission(in *syncnoterpc.GrantPermissionReq) (*syncnoterpc.PermissionResp, error) {
	//想要给予权限，首先先查操作人对于这个note是否有admin以上的权限。
	if in.OperatorId == "" && in.TargetUserId == "" {
		return nil, errors.New("miss param")
	}
	op := sql.NullString{String: in.OperatorId, Valid: true}
	targetid := sql.NullString{String: in.TargetUserId, Valid: true}
	operator, err := l.svcCtx.NotePermissionsModel.FindOneByNoteIdUserId(l.ctx, in.NoteId, op)
	if err != nil {
		return nil, err
	}
	if operator.Role != syncnoterpc.Role_ROLE_OWNER.String() && operator.Role != syncnoterpc.Role_ROLE_ADMIN.String() {
		return nil, errors.New("permission denied")
	}
	res, err := l.svcCtx.NotePermissionsModel.Insert(l.ctx, &model.NotePermissions{
		PermissionId: uuid.New().String(),
		NoteId:       in.NoteId,
		UserId:       targetid,
		Role:         in.Role.String(),
	})
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows != 1 {
		return nil, errors.New("failed to insert permission: no rows affected")
	}
	return &syncnoterpc.PermissionResp{
		Success: true,
		Message: "Success",
	}, nil
}
