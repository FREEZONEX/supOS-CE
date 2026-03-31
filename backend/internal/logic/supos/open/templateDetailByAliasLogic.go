// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/template/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type TemplateDetailByAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询模板详情
func NewTemplateDetailByAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TemplateDetailByAliasLogic {
	return &TemplateDetailByAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TemplateDetailByAliasLogic) TemplateDetailByAlias(alias string) (resp *types.ResultVO, err error) {
	// 调用 UnsTemplateService.DetailByAlias
	result, err := spring.GetBean[*service.UnsTemplateService]().DetailByAlias(l.ctx, &types.WithAlias{Alias: alias})
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.template.query.failed") + ": " + err.Error(),
		}, nil
	}
	return &types.ResultVO{
		Code: result.Code,
		Msg:  result.Msg,
	}, nil
}
