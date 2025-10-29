// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改文件夹或文件名称
func NewUpdateNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNameLogic {
	return &UpdateNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNameLogic) UpdateName(req *types.UpdateNameVo) (resp *types.StringResult, err error) {
	resp, err = spring.GetBean[*service.UnsUpdateService]().UpdateName(l.ctx, req)
	return
}
