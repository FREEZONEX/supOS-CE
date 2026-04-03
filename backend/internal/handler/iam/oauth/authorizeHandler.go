package oauth

import (
	"net/http"

	"backend/internal/common/constants"
	oauthlogic "backend/internal/logic/iam/oauth"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthorizeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OAuthAuthorizeReq
		if err := httpx.Parse(r, &req); err != nil {
			writeSimpleOAuthError(w, http.StatusBadRequest, "invalid_request", "failed to parse authorize request")
			return
		}

		var sessionID string
		if cookie, _ := r.Cookie(constants.AccessTokenKey); cookie != nil {
			sessionID = cookie.Value
		}

		l := oauthlogic.NewAuthorizeLogic(r.Context(), svcCtx)
		resp, err := l.Authorize(&req, sessionID, r.URL.RequestURI())
		if err != nil {
			writeOAuthError(w, err)
			return
		}
		if resp == nil || resp.RedirectURL == "" {
			writeSimpleOAuthError(w, http.StatusInternalServerError, "server_error", "authorize redirect is empty")
			return
		}
		http.Redirect(w, r, resp.RedirectURL, http.StatusFound)
	}
}
