package handler

import (
	"backend/internal/adapters/kong/logic"
	"backend/internal/adapters/kong/vo"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type KongHandler struct {
	kongLogic *logic.KongLogic
}

func NewKongHandler(kongLogic *logic.KongLogic) *KongHandler {
	return &KongHandler{kongLogic: kongLogic}
}

// RouteListHandler 对应 GET /inter-api/supos/kong/routeList
func (h *KongHandler) RouteListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes, err := h.kongLogic.RouteList(r.Context())
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, vo.Success(routes))
	}
}
