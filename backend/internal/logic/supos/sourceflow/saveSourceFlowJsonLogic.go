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

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SaveSourceFlowJsonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Persist Node-RED flow JSON
func NewSaveSourceFlowJsonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSourceFlowJsonLogic {
	return &SaveSourceFlowJsonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveSourceFlowJsonLogic) SaveSourceFlowJson(req *types.SourceFlowSaveReq) error {
	return l.SaveFlowJsonWithType(req, constants.FlowTypeNODERED)
}

// SaveFlowJsonWithType persists flow JSON for the specified flow type.
func (l *SaveSourceFlowJsonLogic) SaveFlowJsonWithType(req *types.SourceFlowSaveReq, flowType string) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	var flowData string
	if len(req.Flows) > 0 {
		data, err := json.Marshal(req.Flows)
		if err != nil {
			return errors.Parameter.WithMsg("nodered.invalid.parameter")
		}
		flowData = string(data)
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	rec, err := LoadFlowByType(l.ctx, repo, flowID, flowType)
	if err != nil {
		return err
	}
	rec.FlowData = flowData
	if strings.TrimSpace(rec.FlowID) != "" {
		rec.FlowStatus = flowcommon.FlowStatusPending
	} else {
		rec.FlowStatus = flowcommon.FlowStatusDraft
	}
	return repo.Update(l.ctx, rec)
}
