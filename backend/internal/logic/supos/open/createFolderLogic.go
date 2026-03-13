// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/constants"
	"backend/internal/common/utils/PathUtil"
	unsService "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建文件夹
func NewCreateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFolderLogic {
	return &CreateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFolderLogic) CreateFolder(req *types.CreateOpenApiFolderDto) (resp *types.ResultVO, err error) {
	// 转换为 CreateTopicDto
	createTopicDto := &types.CreateTopicDto{
		Name:     req.Name,
		PathType: constants.PathTypeDir,
	}

	// 设置 Alias (如果为空则由系统生成)
	if req.Alias != "" {
		createTopicDto.Alias = req.Alias
	} else {
		createTopicDto.Alias = PathUtil.GenerateAlias(req.Name, constants.PathTypeDir)
	}

	// 设置 DisplayName
	if req.DisplayName != "" {
		createTopicDto.DisplayName = base.V2p(req.DisplayName)
	}

	// 设置 ParentAlias
	if req.ParentAlias != "" {
		createTopicDto.ParentAlias = base.V2p(req.ParentAlias)
	}

	// 设置 Description
	if req.Description != "" {
		createTopicDto.Description = base.V2p(req.Description)
	}

	// 设置 TemplateAlias/ModelAlias
	if req.TemplateAlias != "" {
		createTopicDto.ModelAlias = base.V2p(req.TemplateAlias)
	}

	// 设置 ExtendProperties
	if len(req.ExtendProperties) > 0 {
		createTopicDto.Extend = req.ExtendProperties
	}

	// 设置字段定义 - 需要转换 []FieldDefine 到 []*FieldDefine
	if len(req.Definition) > 0 {
		createTopicDto.Fields = make([]*types.FieldDefine, len(req.Definition))
		for i, f := range req.Definition {
			createTopicDto.Fields[i] = &f
		}
	}

	// 调用 UnsAddService
	result := spring.GetBean[*unsService.UnsAddService]().CreateModelInstance(l.ctx, createTopicDto)
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
