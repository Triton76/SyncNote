package handler

import (
	"net/http"

	"SyncNote/rebuild/syncnote/api/internal/logic"
	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 加入团队
func JoinTeamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinTeamRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewJoinTeamLogic(r.Context(), svcCtx)
		resp, err := l.JoinTeam(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
