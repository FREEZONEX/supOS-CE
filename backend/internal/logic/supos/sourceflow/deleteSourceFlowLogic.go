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

type DeleteSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete a source flow by id
func NewDeleteSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSourceFlowLogic {
	return &DeleteSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSourceFlowLogic) DeleteSourceFlow(req *types.SourceFlowDeleteReq) error {
	return l.DeleteFlowWithType(req, constants.FlowTypeNODERED)
}

// DeleteFlowWithType deletes the flow ensuring it belongs to the given template type.
func (l *DeleteSourceFlowLogic) DeleteFlowWithType(req *types.SourceFlowDeleteReq, flowType string) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	template := strings.TrimSpace(flowType)
	if template != "" {
		if _, err := LoadFlowByType(l.ctx, repo, flowID, template); err != nil {
			return err
		}
	}
	if err := repo.ReplaceModels(l.ctx, flowID, nil); err != nil {
		return err
	}
	return repo.Delete(l.ctx, flowID)
}
