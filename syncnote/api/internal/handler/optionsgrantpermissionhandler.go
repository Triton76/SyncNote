package handler

import (
	"net/http"

	"SyncNote/syncnote/api/internal/logic"
	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OptionsGrantPermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserNotesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOptionsGrantPermissionLogic(r.Context(), svcCtx)
		resp, err := l.OptionsGrantPermission(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
