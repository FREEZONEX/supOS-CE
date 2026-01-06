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

type UpdateBaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改模板基本信息
func NewUpdateBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseInfoLogic {
	return &UpdateBaseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBaseInfoLogic) UpdateBaseInfo(req *types.UpdateTemplateBaseInfoReq) (resp *types.BaseResult, err error) {
	return spring.GetBean[*service.UnsTemplateService]().UpdateBaseInfo(l.ctx, req)
}
