// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package label

import (
	"backend/internal/logic/supos/uns/label/service"
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

// 修改标签订阅
func NewUpdateSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubscribeLogic {
	return &UpdateSubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSubscribeLogic) UpdateSubscribe(req *types.UpdateLabelSubscribeReq) (resp *types.BaseResult, err error) {
	sv := spring.GetBean[*service.UnsLabelService]()
	return sv.UpdateSubscribe(l.ctx, req)
}
