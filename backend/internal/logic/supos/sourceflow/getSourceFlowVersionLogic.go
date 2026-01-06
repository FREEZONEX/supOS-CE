// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSourceFlowVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Query current Node-RED flow version
func NewGetSourceFlowVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSourceFlowVersionLogic {
	return &GetSourceFlowVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSourceFlowVersionLogic) GetSourceFlowVersion() (map[string]string, error) {
	client := l.svcCtx.SourceNodeRed
	if client == nil {
		return map[string]string{"rev": ""}, nil
	}
	var out map[string]any
	code, body, errs := client.GetVersionRevV2(l.ctx, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("get source flow version failed: code=%d err=%v body=%s", code, errs, string(body))
		return map[string]string{"rev": ""}, nil
	}
	if rev, ok := out["rev"].(string); ok {
		return map[string]string{"rev": rev}, nil
	}
	return map[string]string{"rev": ""}, nil
}
