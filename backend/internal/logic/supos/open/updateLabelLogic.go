// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/label/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改标签
func NewUpdateLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLabelLogic {
	return &UpdateLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLabelLogic) UpdateLabel(id int64, req *types.UpdateLabelDto) (resp *types.ResultVO, err error) {
	// 转换为 UpdateLabelReq
	updateReq := &types.UpdateLabelReq{
		ID:                 id,
		LabelName:          req.LabelName,
		SubscribeEnable:    base.V2p(req.SubscribeEnable),
		SubscribeFrequency: req.SubscribeFrequency,
	}

	// 调用 UnsLabelService.Update
	result, err := spring.GetBean[*service.UnsLabelService]().Update(l.ctx, updateReq)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.label.update.failed") + ": " + err.Error(),
		}, nil
	}

	if result.Code != 200 {
		return &types.ResultVO{
			Code: result.Code,
			Msg:  result.Msg,
		}, nil
	}

	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
	}, nil
}
