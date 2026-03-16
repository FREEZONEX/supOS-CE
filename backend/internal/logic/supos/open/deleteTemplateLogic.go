// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/template/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除模板
func NewDeleteTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTemplateLogic {
	return &DeleteTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTemplateLogic) DeleteTemplate(id int64) (resp *types.ResultVO, err error) {
	// 调用 UnsTemplateService.Delete
	result, err := spring.GetBean[*service.UnsTemplateService]().Delete(l.ctx, &types.WithID{ID: id})
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.template.delete.failed") + ": " + err.Error(),
		}, nil
	}
	return &types.ResultVO{
		Code: result.Code,
		Msg:  result.Msg,
	}, nil
}
