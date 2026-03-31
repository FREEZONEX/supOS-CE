// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	unsservice "backend/internal/logic/supos/uns/uns/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsTreeByDefinitionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 多条件分页查询树结构
func NewUnsTreeByDefinitionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsTreeByDefinitionsLogic {
	return &UnsTreeByDefinitionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsTreeByDefinitionsLogic) UnsTreeByDefinitions(req *types.UnsTreeCondition) (resp *types.TreePageResult, err error) {
	unsQueryService := spring.GetBean[*unsservice.UnsQueryService]()
	unsResp, err := unsQueryService.LazyTree(l.ctx, req)
	if err != nil {
		return &types.TreePageResult{
			Code: 500,
		}, err
	}

	// Convert UnsTreePageResp to TreePageResult
	resp = &types.TreePageResult{
		PageNo:   unsResp.PageNo,
		PageSize: unsResp.PageSize,
		Total:    unsResp.Total,
		Code:     unsResp.Code,
		Data:     unsResp.Data,
	}
	return
}
