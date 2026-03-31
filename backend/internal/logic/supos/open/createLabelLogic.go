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

type CreateLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewCreateLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLabelLogic {
	return &CreateLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLabelLogic) CreateLabel(req *types.CreateLabelDto) (resp *types.ResultVO, err error) {
	// 转换为 CreateLabelReq
	createReq := &types.CreateLabelReq{
		Name: req.LabelName,
	}

	// 调用 UnsLabelService.Create
	result, err := spring.GetBean[*service.UnsLabelService]().Create(l.ctx, createReq)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.label.create.failed") + ": " + err.Error(),
		}, nil
	}
	return &types.ResultVO{
		Code: result.Code,
		Msg:  result.Msg,
		Data: map[string]interface{}{
			"id": result.Data.ID,
		},
	}, nil
}
