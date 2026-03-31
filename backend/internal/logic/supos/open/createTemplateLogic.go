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

type CreateTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增模板
func NewCreateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTemplateLogic {
	return &CreateTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTemplateLogic) CreateTemplate(req *types.CreateTemplateVo) (resp *types.ResultVO, err error) {
	// 转换为 CreateTemplateReq
	createReq := &types.CreateTemplateReq{
		Name:        req.Name,
		Alias:       req.Alias,
		Fields:      make([]*types.FieldDefine, len(req.Fields)),
		Description: req.Description,
	}

	// 转换 Fields
	for i, f := range req.Fields {
		createReq.Fields[i] = &f
	}

	// 调用 UnsTemplateService.Create
	result, err := spring.GetBean[*service.UnsTemplateService]().Create(l.ctx, createReq)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.template.create.failed") + ": " + err.Error(),
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
		Data: map[string]interface{}{
			"id": result.Id,
		},
	}, nil
}
