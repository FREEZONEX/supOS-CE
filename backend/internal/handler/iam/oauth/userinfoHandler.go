package oauth

import (
	"encoding/json"
	"net/http"

	oauthlogic "backend/internal/logic/iam/oauth"
	"backend/internal/svc"
)

func UserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := oauthlogic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo(bearerToken(r))
		if err != nil {
			writeOAuthError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
