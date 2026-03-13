// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/logic/supos/uns/template/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type TemplatePageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询模板列表
func NewTemplatePageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TemplatePageListLogic {
	return &TemplatePageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TemplatePageListLogic) TemplatePageList(req *types.TemplateQueryVo) (resp *types.TemplatePageResult, err error) {
	// 调用 UnsTemplateService.PageList
	result, err := spring.GetBean[*service.UnsTemplateService]().PageList(l.ctx, req)
	if err != nil {
		return &types.TemplatePageResult{
			Code: 500,
		}, nil
	}

	if result.Code != 200 {
		return &types.TemplatePageResult{
			Code: result.Code,
		}, nil
	}

	// 转换 []*TemplateSearchResult 到 []TemplateSearchResult
	data := make([]types.TemplateSearchResult, len(result.Data))
	for i, item := range result.Data {
		if item != nil {
			data[i] = *item
		}
	}

	return &types.TemplatePageResult{
		Code:     200,
		Data:     data,
		PageNo:   result.PageNo,
		PageSize: result.PageSize,
		Total:    result.Total,
	}, nil
}
