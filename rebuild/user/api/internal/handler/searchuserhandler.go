package handler

import (
	"net/http"

	"SyncNote/rebuild/user/api/internal/logic"
	"SyncNote/rebuild/user/api/internal/svc"
	"SyncNote/rebuild/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSearchUserLogic(r.Context(), svcCtx)
		resp, err := l.SearchUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
