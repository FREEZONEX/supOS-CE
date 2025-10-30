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

type CancelLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件取消标签
func NewCancelLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLabelLogic {
	return &CancelLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLabelLogic) CancelLabel(req *types.CancelLabelReq) (resp *types.BaseResult, err error) {
	sv := spring.GetBean[*service.UnsLabelService]()
	err = sv.CancelLabel(l.ctx, req.UnsId, req.LabelIds)
	if err == nil {
		resp = &types.BaseResult{Code: 200, Msg: "ok"}
	}
	return
}
