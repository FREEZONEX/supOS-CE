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

type MarkEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Mark a  flow by id
func NewMarkEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkEventFlowLogic {
	return &MarkEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkEventFlowLogic) MarkEventFlow(req *types.FlowMarkReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	return sourceflow.NewMarkSourceFlowLogic(l.ctx, l.svcCtx).
		MarkFlowWithType(req, constants.FlowTypeEVENTFLOW)
}
