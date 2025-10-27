package auth

import (
	"net/http"
	"strings"

	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// token
func TokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TokenCallbackReq
		if err := httpx.Parse(r, &req); err != nil || strings.TrimSpace(req.Code) == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		l := authlogic.NewTokenLogic(r.Context(), svcCtx)
		result, err := l.Token(&req)
		if err != nil || result == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if result.Cookie != nil {
			http.SetCookie(w, result.Cookie)
		}
		if result.RedirectURL != "" {
			w.Header().Set("Location", result.RedirectURL)
		}
		w.WriteHeader(http.StatusFound)
	}
}
