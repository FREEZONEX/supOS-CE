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

type DetailByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID查询模板详情
func NewDetailByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByIdLogic {
	return &DetailByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailByIdLogic) DetailById(req *types.WithID) (resp *types.TemplateDetailResp, err error) {
	return spring.GetBean[*service.UnsTemplateService]().DetailById(l.ctx, req)
}
