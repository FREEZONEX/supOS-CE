// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package nodered

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyNodeRedFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 代理 NodeRed /flows 接口
func NewProxyNodeRedFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyNodeRedFlowsLogic {
	return &ProxyNodeRedFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyNodeRedFlowsLogic) ProxyNodeRedFlows() (resp string, err error) {
	// todo: add your logic here and delete this line

	return
}
