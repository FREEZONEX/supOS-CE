// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/template/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改模板
func NewUpdateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateLogic {
	return &UpdateTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTemplateLogic) UpdateTemplate(alias string, req *types.UpdateTemplateDto) (resp *types.ResultVO, err error) {
	// 转换为 UpdateTemplateFieldsAndDescReq
	updateReq := &types.UpdateTemplateFieldsAndDescReq{
		Alias:       alias,
		Fields:      make([]*types.FieldDefine, len(req.Fields)),
		Description: base.V2p(req.Description),
	}

	// 转换 Fields
	for i, f := range req.Fields {
		updateReq.Fields[i] = &f
	}

	// 调用 UnsTemplateService.UpdateFieldsAndDesc
	result, err := spring.GetBean[*service.UnsTemplateService]().UpdateFieldsAndDesc(l.ctx, updateReq)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.template.update.failed") + ": " + err.Error(),
		}, nil
	}
	return &types.ResultVO{
		Code: result.Code,
		Msg:  result.Msg,
	}, nil
}
