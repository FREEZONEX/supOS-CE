package contextx

import "context"

type Subject struct {
	UserID       int64
	UserName     string
	Email        string
	AuthType     string
	APIKeyName   string
	APIKeyPrefix string
	KeyType      string
	ResourceKeys map[string]struct{}
	UIActions    []string
}

type subjectKey struct{}

func WithSubject(ctx context.Context, subject Subject) context.Context {
	return context.WithValue(ctx, subjectKey{}, subject)
}

func SubjectFrom(ctx context.Context) (Subject, bool) {
	subject, ok := ctx.Value(subjectKey{}).(Subject)
	return subject, ok
}

func HasResource(subject Subject, resourceKey string) bool {
	if resourceKey == "" {
		return true
	}
	_, ok := subject.ResourceKeys[resourceKey]
	return ok
}
