package oauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	oauthlogic "backend/internal/logic/iam/oauth"
)

func writeOAuthError(w http.ResponseWriter, err error) {
	var oauthErr *oauthlogic.OAuthError
	if errors.As(err, &oauthErr) && oauthErr != nil {
		status := oauthErr.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":             oauthErr.Code,
			"error_description": oauthErr.Description,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":             "server_error",
		"error_description": "unexpected oauth error",
	})
}

func writeSimpleOAuthError(w http.ResponseWriter, status int, code, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":             code,
		"error_description": description,
	})
}

func bearerToken(r *http.Request) string {
	if r == nil {
		return ""
	}
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return strings.TrimSpace(r.URL.Query().Get("access_token"))
	}
	const prefix = "Bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}
