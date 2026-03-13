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

type DeleteLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewDeleteLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLabelLogic {
	return &DeleteLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLabelLogic) DeleteLabel(id int64) (resp *types.ResultVO, err error) {
	// 调用 UnsLabelService.Delete
	result, err := spring.GetBean[*service.UnsLabelService]().Delete(l.ctx, &types.WithID{ID: id})
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.label.delete.failed") + ": " + err.Error(),
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
