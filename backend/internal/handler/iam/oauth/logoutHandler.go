package oauth

import (
	"net/http"
	"strings"

	oauthlogic "backend/internal/logic/iam/oauth"
	"backend/internal/svc"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			token = strings.TrimSpace(r.FormValue("token"))
		}

		l := oauthlogic.NewLogoutLogic(r.Context(), svcCtx)
		if err := l.Logout(token); err != nil {
			writeOAuthError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
