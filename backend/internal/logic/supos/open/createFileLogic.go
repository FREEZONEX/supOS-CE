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

type CreateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建文件
func NewCreateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFileLogic {
	return &CreateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFileLogic) CreateFile(req *types.CreateOpenApiFileDto) (resp *types.ResultVO, err error) {
	// 转换为 CreateTopicDto
	createTopicDto := &types.CreateTopicDto{
		Name:     req.Name,
		PathType: constants.PathTypeFile,
	}

	// 设置 Alias (如果为空则由系统生成)
	if req.Alias != "" {
		createTopicDto.Alias = req.Alias
	} else {
		createTopicDto.Alias = PathUtil.GenerateAlias(req.Name, constants.PathTypeFile)
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

	// 设置 DataType
	if req.DataType != 0 {
		createTopicDto.DataType = base.V2p(req.DataType)
	}

	// 设置 TemplateAlias/ModelAlias
	if req.TemplateAlias != "" {
		createTopicDto.ModelAlias = base.V2p(req.TemplateAlias)
	}

	// 设置标志位
	createTopicDto.Save2Db = base.V2p(req.Persistence)
	createTopicDto.AddDashBoard = base.V2p(req.DashBoard)
	createTopicDto.AddFlow = base.V2p(req.AddFlow)

	// 设置 Expression
	if req.Expression != "" {
		createTopicDto.Expression = base.V2p(req.Expression)
	}

	// 设置 Frequency
	if req.Frequency != "" {
		createTopicDto.Frequency = req.Frequency
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

	// 设置 Refers (引用对象)
	if len(req.Refers) > 0 {
		referIds := make([]int64, len(req.Refers))
		for i, ref := range req.Refers {
			referIds[i] = ref.Id
		}
		createTopicDto.ReferIds = referIds
	}

	// 设置 ParentDataType
	if req.ParentDataType != 0 {
		createTopicDto.ParentDataType = base.V2p(req.ParentDataType)
	}

	// 设置 LabelNames
	if len(req.LabelNames) > 0 {
		createTopicDto.LabelNames = req.LabelNames
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
