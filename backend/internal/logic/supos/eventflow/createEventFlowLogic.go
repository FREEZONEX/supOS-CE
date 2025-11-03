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
	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create a new event flow
func NewCreateEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEventFlowLogic {
	return &CreateEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEventFlowLogic) CreateEventFlow(req *types.EventFlowCreateReq) (string, error) {
	if req == nil {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	srcReq := &types.SourceFlowCreateReq{
		FlowName:    req.FlowName,
		Description: req.Description,
		Template:    req.Template,
	}
	return sourceflow.NewCreateSourceFlowLogic(l.ctx, l.svcCtx).
		CreateFlowWithType(srcReq, constants.FlowTypeEVENTFLOW)
}
