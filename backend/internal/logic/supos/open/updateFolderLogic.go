// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改文件夹
func NewUpdateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFolderLogic {
	return &UpdateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFolderLogic) UpdateFolder(alias string, req *types.UpdateOpenApiFolderDto) (resp *types.ResultVO, err error) {
	// 转换为 UpdateUnsDto
	updateDto := &types.UpdateUnsDto{
		Alias:       alias,
		Name:        req.Name,
		Description: base.V2p(req.Description),
		ModelAlias:  base.V2p(req.TemplateAlias),
	}

	// 设置 ParentAlias (空字符串时不设置)
	if req.ParentAlias != "" {
		updateDto.ParentAlias = base.V2p(req.ParentAlias)
	}

	// 设置 DisplayName
	if req.DisplayName != "" {
		updateDto.DisplayName = base.V2p(req.DisplayName)
	}

	// 设置 ExtendProperties
	if len(req.ExtendProperties) > 0 {
		updateDto.Extend = req.ExtendProperties
	}

	// 设置字段定义 - 需要转换 []FieldDefine 到 []*FieldDefine
	if len(req.Definition) > 0 {
		updateDto.Fields = make([]*types.FieldDefine, len(req.Definition))
		for i, f := range req.Definition {
			updateDto.Fields[i] = &f
		}
	}

	// 调用 UnsUpdateService
	result, err := spring.GetBean[*service.UnsUpdateService]().UpdateDetail(l.ctx, updateDto)
	if err != nil {
		l.Errorf("修改文件夹失败: %v", err)
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.folder.update.failed") + ": " + err.Error(),
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
