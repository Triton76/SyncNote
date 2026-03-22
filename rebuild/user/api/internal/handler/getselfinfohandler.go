package handler

import (
	"net/http"

	"SyncNote/rebuild/user/api/internal/logic"
	"SyncNote/rebuild/user/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSelfInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetSelfInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetSelfInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
