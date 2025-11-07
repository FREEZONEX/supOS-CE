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

type UnmarkEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete Mark a  flow by id
func NewUnmarkEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarkEventFlowLogic {
	return &UnmarkEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnmarkEventFlowLogic) UnmarkEventFlow(req *types.FlowUNMarkReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	return sourceflow.NewUnmarkSourceFlowLogic(l.ctx, l.svcCtx).
		UnmarkFlowWithType(req, constants.FlowTypeEVENTFLOW)
}
