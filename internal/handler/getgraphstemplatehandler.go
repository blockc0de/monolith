package handler

import (
	"net/http"

	"github.com/blockc0de/monolith/internal/logic"
	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetGraphsTemplateHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGraphsTemplateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetGraphsTemplateLogic(r.Context(), ctx)
		resp, err := l.GetGraphsTemplate(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
