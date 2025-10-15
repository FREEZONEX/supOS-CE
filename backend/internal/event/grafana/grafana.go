package grafana

import (
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrafanaHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewGrafanaHandle(ctx context.Context, svcCtx *svc.ServiceContext) *GrafanaHandle {
	return &GrafanaHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}
