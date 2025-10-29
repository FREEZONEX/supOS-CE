package ctxs

import (
	"context"
	"net/http"
	"time"

	"gitee.com/unitedrhino/share/utils"
	"go.opentelemetry.io/otel/trace"
)

var ContextKeys = []string{UserTokenKey, UserSetTokenKey, UserWorkspaceKey}

func CopyCtx(ctx context.Context) context.Context {
	newCtx := NewUserCtx(ctx)
	newCtx = trace.ContextWithSpanContext(newCtx, trace.SpanContextFromContext(ctx))
	for _, k := range ContextKeys {
		if v := ctx.Value(k); v != nil {
			newCtx = context.WithValue(newCtx, k, v)
		}
	}
	return newCtx
}

func GoNewCtx(ctx context.Context, f func(ctx2 context.Context)) {
	ctx = CopyCtx(ctx)
	go func() {
		defer utils.Recover(ctx)
		f(ctx)
	}()
}

func GetDeadLine(ctx context.Context, defaultDeadLine time.Time) time.Time {
	dead, ok := ctx.Deadline()
	if !ok {
		return defaultDeadLine
	}
	return dead
}

const needRespKey = "result.needResp"

func NeedResp(r *http.Request) *http.Request {
	v := r.Context().Value(needRespKey)
	if v != nil {
		return r
	}
	return r.WithContext(context.WithValue(r.Context(), needRespKey, &http.Response{}))
}
func GetResp(r *http.Request) *http.Response {
	v := r.Context().Value(needRespKey)
	if v == nil {
		return nil
	}
	vv := v.(*http.Response)
	return vv
}
