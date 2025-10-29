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

type UpdateDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改文件夹或文件明细
func NewUpdateDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDetailLogic {
	return &UpdateDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDetailLogic) UpdateDetail(req *types.UpdateUnsDto) (resp *types.StringResult, err error) {
	resp, err = spring.GetBean[*service.UnsUpdateService]().UpdateDetail(l.ctx, req)
	return
}
