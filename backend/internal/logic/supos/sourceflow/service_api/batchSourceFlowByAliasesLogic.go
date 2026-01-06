// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"context"
	"strconv"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchSourceFlowByAliasesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Batch query flows by UNS aliases
func NewBatchSourceFlowByAliasesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSourceFlowByAliasesLogic {
	return &BatchSourceFlowByAliasesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchSourceFlowByAliasesLogic) BatchSourceFlowByAliases(req *types.SourceFlowBatchAliasReq) (resp *types.SourceFlowListResult, err error) {
	if req == nil || len(req.Aliases) == 0 {
		return &types.SourceFlowListResult{
			Code: 0,
			Data: []types.SourceFlowInfo{},
		}, nil
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	flows, err := repo.SelectByAliases(l.ctx, req.Aliases)
	if err != nil {
		return nil, err
	}
	items := make([]types.SourceFlowInfo, 0, len(flows))
	for _, f := range flows {
		if f == nil {
			continue
		}
		items = append(items, types.SourceFlowInfo{
			ID:          strconv.FormatInt(f.ID, 10),
			FlowName:    f.FlowName,
			FlowID:      f.FlowID,
			Description: f.Description,
			FlowStatus:  f.FlowStatus,
			Template:    f.Template,
		})
	}
	return &types.SourceFlowListResult{
		Code: 0,
		Data: items,
	}, nil
}
