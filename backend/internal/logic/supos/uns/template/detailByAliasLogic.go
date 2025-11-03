// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package template

import (
	"backend/internal/logic/supos/uns/template/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailByAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据别名查询模板详情
func NewDetailByAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByAliasLogic {
	return &DetailByAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailByAliasLogic) DetailByAlias(req *types.WithAlias) (resp *types.TemplateDetailResp, err error) {
	return spring.GetBean[*service.UnsTemplateService]().DetailByAlias(l.ctx, req)
}
