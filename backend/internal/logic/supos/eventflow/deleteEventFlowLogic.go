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

type DeleteEventFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete an event flow by id
func NewDeleteEventFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEventFlowLogic {
	return &DeleteEventFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEventFlowLogic) DeleteEventFlow(req *types.EventFlowDeleteReq) error {
	if req == nil {
		return errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.flowId.empty"))
	}
	repo := relationDB.NewNoderedEventFlowRepo(l.ctx)
	if err := repo.ReplaceModels(l.ctx, flowID, nil); err != nil {
		return err
	}
	return repo.Delete(l.ctx, flowID)
}
