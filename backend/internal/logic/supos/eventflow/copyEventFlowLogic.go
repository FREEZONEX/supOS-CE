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

type CopyEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Copy an existing event flow
func NewCopyEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CopyEventFlowLogic {
	return &CopyEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CopyEventFlowLogic) CopyEventFlow(req *types.EventFlowCopyReq) (string, error) {
	if req == nil {
		return "", errors.Parameter.WithMsg("error.sys.parameterError")
	}
	srcReq := &types.SourceFlowCopyReq{
		SourceID:    req.SourceID,
		FlowName:    req.FlowName,
		Description: req.Description,
		Template:    constants.FlowTypeEVENTFLOW,
	}
	return sourceflow.NewCopySourceFlowLogic(l.ctx, l.svcCtx).
		CopyFlowWithType(srcReq, constants.FlowTypeEVENTFLOW, l.svcCtx.EventNodeRed)
}
