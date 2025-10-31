// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/logic/supos/flowcommon"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/i18ns"
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
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	srcID, err := strconv.ParseInt(strings.TrimSpace(req.SourceID), 10, 64)
	if err != nil || srcID <= 0 {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	repo := relationDB.NewNoderedEventFlowRepo(l.ctx)
	factory := func() *relationDB.NoderedEventFlow {
		return &relationDB.NoderedEventFlow{}
	}
	input := flowcommon.FlowCopyInput{
		FlowName:    name,
		Description: strings.TrimSpace(req.Description),
		Template:    strings.TrimSpace(req.Template),
	}
	record, err := flowcommon.CopyFlow(l.ctx, l.svcCtx, repo, factory, srcID, input, l.svcCtx.EventNodeRed)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(record.GetID(), 10), nil
}
