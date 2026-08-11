package permission

import (
	"context"
	"net/http"

	"backend/internal/contextx"
	"backend/internal/repo"
)

type Evaluator struct {
	iam *repo.IAMRepo
}

func New(ctx context.Context) *Evaluator {
	return &Evaluator{iam: repo.NewIAMRepo(ctx)}
}

func (e *Evaluator) Allow(ctx context.Context, subject contextx.Subject, actionType, method, path string) (bool, string, error) {
	if e == nil || e.iam == nil {
		return false, "", nil
	}
	matches, err := e.iam.MatchAction(ctx, actionType, method, path)
	if err != nil {
		return false, "", err
	}
	if len(matches) == 0 {
		return false, "", nil
	}
	for _, match := range matches {
		if contextx.HasResource(subject, match.ResourceKey) {
			return true, match.ResourceKey, nil
		}
	}
	return false, matches[0].ResourceKey, nil
}

func (e *Evaluator) IsAdmin(ctx context.Context, subject contextx.Subject) (bool, error) {
	if e == nil || e.iam == nil {
		return false, nil
	}
	return e.iam.IsAdmin(ctx, subject.UserID)
}

func IsReadMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}
