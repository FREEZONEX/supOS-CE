// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"encoding/json"
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

type DeploySourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Deploy a source flow
func NewDeploySourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeploySourceFlowLogic {
	return &DeploySourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeploySourceFlowLogic) DeploySourceFlow(req *types.SourceFlowDeployReq) (*types.SourceFlowDeployResult, error) {
	return l.DeployFlowWithType(req, constants.FlowTypeNODERED, l.svcCtx.SourceNodeRed)
}

// DeployFlowWithType deploys a flow definition using the provided Node-RED client.
func (l *DeploySourceFlowLogic) DeployFlowWithType(req *types.SourceFlowDeployReq, flowType string, client *noderedclient.Client) (*types.SourceFlowDeployResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("error.sys.parameterError")
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return nil, errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	var override string
	if len(req.Flows) > 0 {
		data, err := json.Marshal(req.Flows)
		if err != nil {
			return nil, errors.Parameter.WithMsg("nodered.invalid.parameter")
		}
		override = string(data)
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	if _, err := LoadFlowByType(l.ctx, repo, flowID, flowType); err != nil {
		return nil, err
	}
	newFlowID, err := flowcommon.DeployFlow(l.ctx, repo, flowID, override, client, flowcommon.ExtractAliases)
	if err != nil {
		return nil, err
	}
	return &types.SourceFlowDeployResult{
		FlowID: newFlowID,
	}, nil
}
