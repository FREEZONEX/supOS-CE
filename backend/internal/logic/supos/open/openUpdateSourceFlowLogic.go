// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenUpdateSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新 Flow
func NewOpenUpdateSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenUpdateSourceFlowLogic {
	return &OpenUpdateSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenUpdateSourceFlowLogic) OpenUpdateSourceFlow(req *types.OpenSourceFlowUpdateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
