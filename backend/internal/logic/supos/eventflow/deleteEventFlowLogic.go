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

type DeleteEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete an event flow by id
func NewDeleteEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEventFlowLogic {
	return &DeleteEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEventFlowLogic) DeleteEventFlow(req *types.EventFlowDeleteReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	srcReq := &types.SourceFlowDeleteReq{ID: req.ID}
	return sourceflow.NewDeleteSourceFlowLogic(l.ctx, l.svcCtx).
		DeleteFlowWithType(srcReq, constants.FlowTypeEVENTFLOW)
}
