// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update flow metadata
func NewUpdateSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSourceFlowLogic {
	return &UpdateSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSourceFlowLogic) UpdateSourceFlow(req *types.SourceFlowUpdateReq) error {
	return l.UpdateFlowWithType(req, constants.FlowTypeNODERED)
}

// UpdateFlowWithType updates the flow metadata for the specified template.
func (l *UpdateSourceFlowLogic) UpdateFlowWithType(req *types.SourceFlowUpdateReq, flowType string) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	template := strings.TrimSpace(flowType)
	rec, err := LoadFlowByType(l.ctx, repo, flowID, template)
	if err != nil {
		return err
	}
	if !strings.EqualFold(rec.FlowName, name) {
		filter := relationDB.NoderedSourceFlowFilter{
			Name:     name,
			Template: template,
		}
		if exist, err := repo.FindOneByFilter(l.ctx, filter); err == nil && exist != nil && exist.ID != flowID {
			return errors.Duplicate.WithMsg("nodered.flowName.has.used")
		} else if err != nil && !errors.Cmp(err, errors.NotFind) {
			return err
		}
		rec.FlowName = name
	}
	rec.Description = strings.TrimSpace(req.Description)
	return repo.Update(l.ctx, rec)
}
