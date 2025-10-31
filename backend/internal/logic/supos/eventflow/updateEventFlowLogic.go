// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update event flow metadata
func NewUpdateEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventFlowLogic {
	return &UpdateEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEventFlowLogic) UpdateEventFlow(req *types.EventFlowUpdateReq) error {
	if req == nil {
		return errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.flowId.empty"))
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		return errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	repo := relationDB.NewNoderedEventFlowRepo(l.ctx)
	rec, err := repo.FindOne(l.ctx, flowID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(rec.FlowName, name) {
		filter := relationDB.NoderedEventFlowFilter{
			Name: name,
			// FlowType: eventFlowType,
		}
		if exist, err := repo.FindOneByFilter(l.ctx, filter); err == nil && exist != nil && exist.ID != flowID {
			return errors.Duplicate.WithMsg(i18ns.LocalizeMsg("nodered.flowName.has.used"))
		} else if err != nil && !errors.Cmp(err, errors.NotFind) {
			return err
		}
		rec.FlowName = name
	}
	rec.Description = strings.TrimSpace(req.Description)
	return repo.Update(l.ctx, rec)
}
