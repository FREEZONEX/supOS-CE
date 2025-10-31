// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
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
	if req == nil {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("error.sys.parameterError"))
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	filter := relationDB.NoderedSourceFlowFilter{
		Name: name,
		// FlowType: sourceFlowType,
	}
	if exist, err := repo.FindOneByFilter(l.ctx, filter); err == nil && exist != nil {
		return "", errors.Duplicate.WithMsg(i18ns.LocalizeMsg("nodered.flowName.has.used"))
	} else if err != nil && !errors.Cmp(err, errors.NotFind) {
		return "", err
	}
	rec := &relationDB.NoderedSourceFlow{
		ID:          l.svcCtx.SnowFlake.GetSnowflakeId(),
		FlowName:    name,
		Description: strings.TrimSpace(req.Description),
		Template:    strings.TrimSpace(req.Template),
		// FlowType:    sourceFlowType,
		FlowStatus: flowcommon.FlowStatusDraft,
	}
	if err := repo.Insert(l.ctx, rec); err != nil {
		return "", err
	}

	return strconv.FormatInt(rec.ID, 10), nil
}
