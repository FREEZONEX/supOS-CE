// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenGetSourceFlowDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询 Flow Data
func NewOpenGetSourceFlowDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenGetSourceFlowDataLogic {
	return &OpenGetSourceFlowDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenGetSourceFlowDataLogic) OpenGetSourceFlowData(req *types.OpenSourceFlowDataQuery) (resp *types.OpenSourceFlowDataResp, err error) {
	// todo: add your logic here and delete this line

	return
}
