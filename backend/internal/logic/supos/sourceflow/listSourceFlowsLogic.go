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
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListSourceFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List source flows with optional fuzzy search
func NewListSourceFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSourceFlowsLogic {
	return &ListSourceFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSourceFlowsLogic) ListSourceFlows(req *types.SourceFlowListQuery) (*types.SourceFlowPageResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("request is nil")
	}
	pageNo := req.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	page := &stores.PageInfo{
		Page: pageNo,
		Size: pageSize,
	}
	filter := relationDB.NoderedSourceFlowFilter{
		NameLike: strings.TrimSpace(req.Keyword),
		// FlowType: sourceFlowType,
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	list, err := repo.FindByFilter(l.ctx, filter, page)
	if err != nil {
		return nil, err
	}
	total, err := repo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]types.SourceFlowInfo, 0, len(list))
	for _, v := range list {
		if v == nil {
			continue
		}
		items = append(items, types.SourceFlowInfo{
			ID:          strconv.FormatInt(v.ID, 10),
			FlowName:    v.FlowName,
			FlowID:      v.FlowID,
			Description: v.Description,
			FlowStatus:  v.FlowStatus,
			Template:    v.Template,
		})
	}
	return &types.SourceFlowPageResult{
		Code:     0,
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
		Data:     items,
	}, nil
}
