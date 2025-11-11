// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/logic/supos/flowcommon"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/clients/nodered/templates"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSourceMockFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create a mock flow from UNS path
func NewCreateSourceMockFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSourceMockFlowLogic {
	return &CreateSourceMockFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSourceMockFlowLogic) CreateSourceMockFlow(req *types.SourceFlowMockReq) (resp *types.SourceFlowInfo, err error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("error.sys.parameterError")
	}
	alias := strings.TrimSpace(req.UnsAlias)
	path := strings.TrimSpace(req.Path)
	if alias == "" || path == "" {
		return nil, errors.Parameter.WithMsg("error.sys.parameterError")
	}
	tpl, err := templates.Load("mock_metrics.json.tpl")
	if err != nil {
		l.Errorf("load mock template failed: %v", err)
		return nil, errors.System.AddDetail(err)
	}
	rendered := templates.RenderDollar(tpl, map[string]string{
		"uns_path":         path,
		"alias_path_topic": path,
		"payload":          "",
		"disabled":         "false",
		"clientid":         alias,
	}, flowcommon.GenerateNodeID)

	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	flowName, _, err := repo.FindAvailableFlowName(l.ctx, alias, constants.FlowTypeNODERED)
	if err != nil {
		return nil, err
	}
	rec := &relationDB.NoderedSourceFlow{
		ID:          common.NextId(),
		FlowName:    flowName,
		Description: fmt.Sprintf("auto mock for %s", alias),
		Template:    constants.FlowTypeNODERED,
		// FlowType:    sourceFlowType,
		FlowStatus: flowcommon.FlowStatusDraft,
		FlowData:   rendered,
	}
	if err := repo.Insert(l.ctx, rec); err != nil {
		return nil, err
	}
	if _, err := flowcommon.DeployFlow(l.ctx, repo, rec.ID, rendered, l.svcCtx.SourceNodeRed, flowcommon.ExtractAliases); err != nil {
		return nil, err
	}
	updated, err := repo.FindOne(l.ctx, rec.ID)
	if err != nil {
		return nil, err
	}
	return &types.SourceFlowInfo{
		ID:          strconv.FormatInt(updated.ID, 10),
		FlowName:    updated.FlowName,
		FlowID:      updated.FlowID,
		Description: updated.Description,
		FlowStatus:  updated.FlowStatus,
		Template:    updated.Template,
	}, nil
}
