// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenCreateSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建 Flow
func NewOpenCreateSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenCreateSourceFlowLogic {
	return &OpenCreateSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenCreateSourceFlowLogic) OpenCreateSourceFlow(req *types.OpenSourceFlowCreateReq) (resp string, err error) {
	// todo: add your logic here and delete this line

	return
}
