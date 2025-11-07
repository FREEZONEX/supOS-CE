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

type UpdateEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update event flow metadata
func NewUpdateEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventFlowLogic {
	return &UpdateEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEventFlowLogic) UpdateEventFlow(req *types.EventFlowUpdateReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	srcReq := &types.SourceFlowUpdateReq{
		ID:          req.ID,
		FlowName:    req.FlowName,
		Description: req.Description,
	}
	return sourceflow.NewUpdateSourceFlowLogic(l.ctx, l.svcCtx).
		UpdateFlowWithType(srcReq, constants.FlowTypeEVENTFLOW)
}
