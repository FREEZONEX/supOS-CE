package oauth

import (
	"encoding/json"
	"net/http"
	"strings"

	oauthlogic "backend/internal/logic/iam/oauth"
	"backend/internal/svc"
	"backend/internal/types"
)

func TokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			writeSimpleOAuthError(w, http.StatusBadRequest, "invalid_request", "failed to parse token request")
			return
		}

		req := types.OAuthTokenReq{
			GrantType:    strings.TrimSpace(r.FormValue("grant_type")),
			Code:         strings.TrimSpace(r.FormValue("code")),
			RedirectURI:  strings.TrimSpace(r.FormValue("redirect_uri")),
			ClientID:     strings.TrimSpace(r.FormValue("client_id")),
			ClientSecret: strings.TrimSpace(r.FormValue("client_secret")),
		}
		if clientID, clientSecret, ok := r.BasicAuth(); ok {
			if req.ClientID == "" {
				req.ClientID = strings.TrimSpace(clientID)
			}
			if req.ClientSecret == "" {
				req.ClientSecret = strings.TrimSpace(clientSecret)
			}
		}

		l := oauthlogic.NewTokenLogic(r.Context(), svcCtx)
		resp, err := l.Token(&req)
		if err != nil {
			writeOAuthError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
