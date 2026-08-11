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

type ReadyzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadyzLogic) Readyz() (resp *types.Envelope, err error) {
	if err := l.svcCtx.App.Ready(l.ctx); err != nil {
		return nil, respx.NewHTTPError(503, err.Error())
	}
	return respx.Envelope(l.svcCtx.App.Readiness(l.ctx)), nil
}
