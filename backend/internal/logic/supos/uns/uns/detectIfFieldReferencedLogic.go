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

type DetectIfFieldReferencedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 预先判断是否有属性关联
func NewDetectIfFieldReferencedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetectIfFieldReferencedLogic {
	return &DetectIfFieldReferencedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetectIfFieldReferencedLogic) DetectIfFieldReferenced(req *types.UpdateModeRequestVo) (resp *types.ResultVO, err error) {
	resp, err = spring.GetBean[*service.UnsQueryService]().DetectIfFieldReferenced(l.ctx, req)
	return
}
