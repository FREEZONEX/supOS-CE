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

type UnmarkSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete Mark a source flow by id
func NewUnmarkSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarkSourceFlowLogic {
	return &UnmarkSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnmarkSourceFlowLogic) UnmarkSourceFlow(req *types.FlowUNMarkReq) error {
	return l.UnmarkFlowWithType(req, constants.FlowTypeNODERED)
}

// UnmarkFlowWithType removes the pin for the specified flow type.
func (l *UnmarkSourceFlowLogic) UnmarkFlowWithType(req *types.FlowUNMarkReq, flowType string) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	idStr := strings.TrimSpace(req.ID)
	flowID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || flowID <= 0 {
		return errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	userID := resolveUserID(l.ctx)
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	if _, err := LoadFlowByType(l.ctx, repo, flowID, flowType); err != nil {
		return err
	}
	topRepo := relationDB.NewNoderedFlowTopRepo(l.ctx)
	if err := topRepo.DeleteByUser(l.ctx, flowID, userID); err != nil {
		return err
	}
	return nil
}
