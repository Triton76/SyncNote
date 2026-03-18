// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"SyncNote/syncnote/api/internal/logic"
	"SyncNote/syncnote/api/internal/svc"
	"SyncNote/syncnote/api/internal/types"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetNoteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NoteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetNoteLogic(r.Context(), svcCtx)
		resp, err := l.GetNote(&req)
		if err != nil {
			if errors.Is(err, logic.ErrForbiddenNoteAccess) {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusForbidden, map[string]any{
					"code":    http.StatusForbidden,
					"message": "forbidden",
				})
				return
			} //然后在这里新增是为了让前端测试更清晰，映射为403.
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
