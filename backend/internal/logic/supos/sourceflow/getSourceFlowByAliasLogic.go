// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSourceFlowByAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Fetch flow information by UNS alias
func NewGetSourceFlowByAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSourceFlowByAliasLogic {
	return &GetSourceFlowByAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSourceFlowByAliasLogic) GetSourceFlowByAlias(req *types.SourceFlowAliasQuery) (resp *types.SourceFlowInfo, err error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("error.sys.parameterError")
	}
	alias := strings.TrimSpace(req.Alias)
	if alias == "" {
		return nil, errors.Parameter.WithMsg("error.sys.parameterError")
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	flow, err := repo.FindLatestByAlias(l.ctx, alias)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, nil
	}
	return &types.SourceFlowInfo{
		ID:          strconv.FormatInt(flow.ID, 10),
		FlowName:    flow.FlowName,
		FlowID:      flow.FlowID,
		Description: flow.Description,
		FlowStatus:  flow.FlowStatus,
		Template:    flow.Template,
	}, nil
}
