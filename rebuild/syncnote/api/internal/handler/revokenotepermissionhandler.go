package handler

import (
	"net/http"

	"SyncNote/rebuild/syncnote/api/internal/logic"
	"SyncNote/rebuild/syncnote/api/internal/svc"
	"SyncNote/rebuild/syncnote/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 撤销笔记权限
func RevokeNotePermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RevokeNotePermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRevokeNotePermissionLogic(r.Context(), svcCtx)
		resp, err := l.RevokeNotePermission(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
