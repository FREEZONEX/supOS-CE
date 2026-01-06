// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventFlowVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Query current event flow version
func NewGetEventFlowVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventFlowVersionLogic {
	return &GetEventFlowVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventFlowVersionLogic) GetEventFlowVersion() (map[string]string, error) {
	client := l.svcCtx.EventNodeRed
	if client == nil {
		return map[string]string{"rev": ""}, nil
	}
	var out map[string]any
	code, body, errs := client.GetVersionRevV2(l.ctx, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("get event flow version failed: code=%d err=%v body=%s", code, errs, string(body))
		return map[string]string{"rev": ""}, nil
	}
	if rev, ok := out["rev"].(string); ok {
		return map[string]string{"rev": rev}, nil
	}
	return map[string]string{"rev": ""}, nil
}
