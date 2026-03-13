// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/label/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

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

func (l *CancelLabelLogic) CancelLabel(alias string, req *types.StringArrayRequest) (resp *types.ResultVO, err error) {
	// 调用 UnsLabelService.CancelLabelByNames
	err = spring.GetBean[*service.UnsLabelService]().CancelLabelByNames(l.ctx, alias, req.Items)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.label.cancel.failed") + ": " + err.Error(),
		}, nil
	}

	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
	}, nil
}
