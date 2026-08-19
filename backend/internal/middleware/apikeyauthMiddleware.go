// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package middleware

import (
	"net/http"

	"backend/internal/contextx"
	"backend/internal/domain/audit"
	authdomain "backend/internal/domain/auth"
	"backend/internal/permission"
)

type ApiKeyAuthMiddleware struct {
	auth       *authdomain.Service
	permission *permission.Evaluator
	audit      *audit.Service
}

func NewApiKeyAuthMiddleware(auth *authdomain.Service, evaluator *permission.Evaluator, auditSvc *audit.Service) *ApiKeyAuthMiddleware {
	return &ApiKeyAuthMiddleware{auth: auth, permission: evaluator, audit: auditSvc}
}

func (m *ApiKeyAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject, err := m.auth.AuthenticateAPIKey(r.Context(), authdomain.APIKeyFromRequest(r))
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		allowed, _, err := m.permission.Allow(r.Context(), subject, "openapi", r.Method, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		ctx := contextx.WithSubject(r.Context(), subject)
		if m.audit != nil {
			m.audit.BindSubject(ctx, subject)
		}
		next(w, r.WithContext(ctx))
	}
}
