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

type UpdateSubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改模板订阅
func NewUpdateSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubscribeLogic {
	return &UpdateSubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSubscribeLogic) UpdateSubscribe(req *types.UpdateTemplateSubscribeReq) (resp *types.BaseResult, err error) {
	return spring.GetBean[*service.UnsTemplateService]().UpdateSubscribe(l.ctx, req)
}
