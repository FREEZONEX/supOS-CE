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

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除模板
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.WithID) (resp *types.BaseResult, err error) {
	return spring.GetBean[*service.UnsTemplateService]().Delete(l.ctx, req)
}
