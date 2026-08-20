// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package system

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthzLogic {
	return &HealthzLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthzLogic) Healthz() (resp *types.Envelope, err error) {
	return respx.Envelope(map[string]any{
		"status": "ok",
		"name":   l.svcCtx.Config.Name,
	}), nil
}
