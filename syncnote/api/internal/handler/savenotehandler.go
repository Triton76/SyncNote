// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"SyncNote/api/internal/logic"
	"SyncNote/api/internal/svc"
	"SyncNote/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SaveNoteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveNoteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSaveNoteLogic(r.Context(), svcCtx)
		resp, err := l.SaveNote(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
