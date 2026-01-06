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

type UpdateFieldsAndDescLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改模板字段（只支持删除和新增）和描述
func NewUpdateFieldsAndDescLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFieldsAndDescLogic {
	return &UpdateFieldsAndDescLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFieldsAndDescLogic) UpdateFieldsAndDesc(req *types.UpdateTemplateFieldsAndDescReq) (resp *types.BaseResult, err error) {
	return spring.GetBean[*service.UnsTemplateService]().UpdateFieldsAndDesc(l.ctx, req)
}
