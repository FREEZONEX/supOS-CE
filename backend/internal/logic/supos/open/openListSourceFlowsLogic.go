// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenListSourceFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询 Flow 列表
func NewOpenListSourceFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenListSourceFlowsLogic {
	return &OpenListSourceFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenListSourceFlowsLogic) OpenListSourceFlows(req *types.OpenSourceFlowListQuery) (resp *types.OpenSourceFlowPageResult, err error) {
	// todo: add your logic here and delete this line

	return
}
