// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/logic/supos/flowcommon"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	noderedclient "backend/share/clients/nodered"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CopySourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Copy an existing source flow
func NewCopySourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CopySourceFlowLogic {
	return &CopySourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CopySourceFlowLogic) CopySourceFlow(req *types.SourceFlowCopyReq) (string, error) {
	return l.CopyFlowWithType(req, constants.FlowTypeNODERED, l.svcCtx.SourceNodeRed)
}

// CopyFlowWithType duplicates a flow for the specified template/flow type.
func (l *CopySourceFlowLogic) CopyFlowWithType(req *types.SourceFlowCopyReq, flowType string, client *noderedclient.Client) (string, error) {
	if req == nil {
		return "", errors.Parameter.WithMsg("error.sys.parameterError")
	}
	srcID, err := strconv.ParseInt(strings.TrimSpace(req.SourceID), 10, 64)
	if err != nil || srcID <= 0 {
		return "", errors.Parameter.WithMsg("error.sys.parameterError")
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		return "", errors.Parameter.WithMsg("error.sys.parameterError")
	}
	template := strings.TrimSpace(req.Template)
	if template == "" {
		template = strings.TrimSpace(flowType)
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	filter := relationDB.NoderedSourceFlowFilter{
		Name:     name,
		Template: template,
	}
	if exist, err := repo.FindOneByFilter(l.ctx, filter); err == nil && exist != nil {
		return "", errors.Duplicate.WithMsg("nodered.flowName.has.used")
	} else if err != nil && !errors.Cmp(err, errors.NotFind) {
		return "", err
	}
	input := flowcommon.FlowCopyInput{
		FlowName:    name,
		Description: strings.TrimSpace(req.Description),
		Template:    template,
	}
	record, err := flowcommon.CopyFlow(l.ctx, l.svcCtx, repo, srcID, input, client)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(record.ID, 10), nil
}
