// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"context"

	"backend/internal/common/constants"
	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListEventFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List event flows with optional fuzzy search
func NewListEventFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEventFlowsLogic {
	return &ListEventFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEventFlowsLogic) ListEventFlows(req *types.EventFlowListQuery) (*types.SourceFlowPageResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("request is nil")
	}
	srcReq := &types.SourceFlowListQuery{
		Keyword:   req.Keyword,
		OrderCode: req.OrderCode,
		IsAsc:     req.IsAsc,
		PageNo:    req.PageNo,
		PageSize:  req.PageSize,
	}
	srcResp, err := sourceflow.NewListSourceFlowsLogic(l.ctx, l.svcCtx).
		ListFlowsWithType(srcReq, constants.FlowTypeEVENTFLOW)
	if err != nil {
		return nil, err
	}
	return srcResp, nil
}
