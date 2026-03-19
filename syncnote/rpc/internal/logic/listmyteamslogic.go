package logic

import (
	"context"
	"errors"

	"SyncNote/syncnote/rpc/internal/svc"
	"SyncNote/syncnote/rpc/pb/syncnoterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyTeamsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyTeamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyTeamsLogic {
	return &ListMyTeamsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Team Membership ---
func (l *ListMyTeamsLogic) ListMyTeams(in *syncnoterpc.ListMyTeamsReq) (*syncnoterpc.ListMyTeamsResp, error) {
	if in == nil || in.UserId == "" {
		return nil, errors.New("userId is required")
	}

	type teamRow struct {
		TeamID   string `db:"team_id"`
		TeamName string `db:"team_name"`
		Role     string `db:"role"`
		Status   string `db:"status"`
		JoinedAt int64  `db:"joined_at"`
	}

	var rows []teamRow
	query := `
		select tm.team_id, coalesce(t.name, '') as team_name, tm.role, tm.status, tm.joined_at
		from team_members tm
		left join teams t on t.team_id = tm.team_id
		where tm.user_id = ? and tm.status = 'active'
		order by tm.joined_at desc`
	if err := l.svcCtx.Conn.QueryRowsCtx(l.ctx, &rows, query, in.UserId); err != nil {
		return nil, err
	}

	teams := make([]*syncnoterpc.TeamInfo, 0, len(rows))
	for _, row := range rows {
		teams = append(teams, &syncnoterpc.TeamInfo{
			TeamId:   row.TeamID,
			TeamName: row.TeamName,
			Role:     row.Role,
			Status:   row.Status,
			JoinedAt: row.JoinedAt,
		})
	}

	return &syncnoterpc.ListMyTeamsResp{Teams: teams}, nil
}
