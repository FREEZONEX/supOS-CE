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

type DeployEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Deploy an event flow
func NewDeployEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeployEventFlowLogic {
	return &DeployEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeployEventFlowLogic) DeployEventFlow(req *types.EventFlowDeployReq) (map[string]string, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("error.sys.parameterError")
	}
	srcReq := &types.SourceFlowDeployReq{
		ID:    req.ID,
		Flows: req.Flows,
	}
	result, err := sourceflow.NewDeploySourceFlowLogic(l.ctx, l.svcCtx).
		DeployFlowWithType(srcReq, constants.FlowTypeEVENTFLOW, l.svcCtx.EventNodeRed)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"flowId": result.FlowID,
	}, nil
}
