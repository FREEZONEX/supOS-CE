// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/logic/supos/flowcommon"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSourceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create a new source flow
func NewCreateSourceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSourceFlowLogic {
	return &CreateSourceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSourceFlowLogic) CreateSourceFlow(req *types.SourceFlowCreateReq) (string, error) {
	return l.CreateFlowWithType(req, constants.FlowTypeNODERED)
}

// CreateFlowWithType creates a flow for the given template/flow type.
func (l *CreateSourceFlowLogic) CreateFlowWithType(req *types.SourceFlowCreateReq, flowType string) (string, error) {
	return l.createFlow(req, strings.TrimSpace(flowType))
}

func (l *CreateSourceFlowLogic) createFlow(req *types.SourceFlowCreateReq, flowType string) (string, error) {
	if req == nil {
		return "", errors.Parameter.WithMsg("error.sys.parameterError")
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		return "", errors.Parameter.WithMsg("error.sys.parameterError")
	}
	template := flowType
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
	rec := &relationDB.NoderedSourceFlow{
		ID:          common.NextId(),
		FlowName:    name,
		Description: strings.TrimSpace(req.Description),
		Template:    template,
		FlowStatus:  flowcommon.FlowStatusDraft,
	}
	userVo := resolveUser(l.ctx)
	if userVo != nil {
		rec.Creator = userVo.PreferredUsername
	}
	if err := repo.Insert(l.ctx, rec); err != nil {
		return "", err
	}

	return strconv.FormatInt(rec.ID, 10), nil
}
