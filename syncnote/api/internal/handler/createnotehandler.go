// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"SyncNote/syncnote/api/internal/logic"
	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateNoteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateNoteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateNoteLogic(r.Context(), svcCtx)
		resp, err := l.CreateNote(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
