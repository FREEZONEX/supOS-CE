// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"encoding/json"
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
	if req == nil {
		return nil, errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return nil, errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.flowId.empty"))
	}
	var override string
	if len(req.Flows) > 0 {
		data, err := json.Marshal(req.Flows)
		if err != nil {
			return nil, errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.invalid.parameter"))
		}
		override = string(data)
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	newFlowID, err := flowcommon.DeployFlow(l.ctx, repo, flowID, override, l.svcCtx.SourceNodeRed, flowcommon.ExtractAliases)
	if err != nil {
		return nil, err
	}
	return &types.SourceFlowDeployResult{
		FlowID: newFlowID,
	}, nil
}
