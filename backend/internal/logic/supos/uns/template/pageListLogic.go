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

type PageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询模板列表
func NewPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageListLogic {
	return &PageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageListLogic) PageList(req *types.TemplateQueryVo) (resp *types.TemplatePageResp, err error) {
	return spring.GetBean[*service.UnsTemplateService]().PageList(l.ctx, req)
}
