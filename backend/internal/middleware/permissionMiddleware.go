// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package middleware

import (
	"net/http"

	"backend/internal/contextx"
	apphttpx "backend/internal/httpx"
	"backend/internal/permission"
)

type PermissionMiddleware struct {
	evaluator *permission.Evaluator
}

func NewPermissionMiddleware(evaluator *permission.Evaluator) *PermissionMiddleware {
	return &PermissionMiddleware{evaluator: evaluator}
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject, ok := contextx.SubjectFrom(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		allowed, _, err := m.evaluator.Allow(r.Context(), subject, "api", r.Method, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !allowed {
			apphttpx.WriteError(w, apphttpx.NewHTTPError(http.StatusForbidden, "forbidden"))
			return
		}
		next(w, r)
	}
}
