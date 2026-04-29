// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenDeploySourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 部署 Flow
func NewOpenDeploySourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenDeploySourceFlowLogic {
	return &OpenDeploySourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenDeploySourceFlowLogic) OpenDeploySourceFlow(req *types.OpenSourceFlowDeployReq) (resp *types.OpenSourceFlowDeployResult, err error) {
	// todo: add your logic here and delete this line

	return
}
