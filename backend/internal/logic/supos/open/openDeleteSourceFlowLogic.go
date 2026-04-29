// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenDeleteSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除 Flow
func NewOpenDeleteSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenDeleteSourceFlowLogic {
	return &OpenDeleteSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenDeleteSourceFlowLogic) OpenDeleteSourceFlow(req *types.OpenSourceFlowDeleteReq) error {
	// todo: add your logic here and delete this line

	return nil
}
