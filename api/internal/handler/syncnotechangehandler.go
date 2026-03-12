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

func SyncNoteChangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NoteChangeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSyncNoteChangeLogic(r.Context(), svcCtx)
		resp, err := l.SyncNoteChange(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
