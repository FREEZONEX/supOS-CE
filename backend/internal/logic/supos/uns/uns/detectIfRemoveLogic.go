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

type DetectIfRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除前预先判断是否有被引用对象
func NewDetectIfRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetectIfRemoveLogic {
	return &DetectIfRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetectIfRemoveLogic) DetectIfRemove(req *types.DetectRemoveReq) (resp *types.RemoveResult, err error) {
	resp, err = spring.GetBean[*service.UnsRemoveService]().DetectIfRemove(l.ctx, req)
	return
}
