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

func (l *ListEventFlowsLogic) ListEventFlows(req *types.EventFlowListQuery) (*types.EventFlowPageResult, error) {
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
	items := make([]types.EventFlowInfo, 0, len(srcResp.Data))
	for _, v := range srcResp.Data {
		items = append(items, types.EventFlowInfo{
			ID:          v.ID,
			FlowName:    v.FlowName,
			FlowID:      v.FlowID,
			Description: v.Description,
			FlowStatus:  v.FlowStatus,
			Template:    v.Template,
			Mark:        v.Mark,
		})
	}
	return &types.EventFlowPageResult{
		Code:     srcResp.Code,
		PageNo:   srcResp.PageNo,
		PageSize: srcResp.PageSize,
		Total:    srcResp.Total,
		Data:     items,
	}, nil
}
