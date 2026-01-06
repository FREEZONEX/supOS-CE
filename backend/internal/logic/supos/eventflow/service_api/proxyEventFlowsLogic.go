// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyEventFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy Node-RED /flows endpoint using cookie scoped id
func NewProxyEventFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyEventFlowsLogic {
	return &ProxyEventFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyEventFlowsLogic) ProxyEventFlows() (resp string, err error) {
	// todo: add your logic here and delete this line

	return
}
