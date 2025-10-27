package UnsConverter

import (
	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	"backend/internal/common/enums"
	"backend/internal/common/utils/FieldUtils"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"strconv"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"github.com/jinzhu/copier"
)

// GetFrequencySeconds 获取频率对应的秒数
func GetFrequencySeconds(frequency string) *int64 {
	nano, ok := enums.TimeUnitsParseToNanoSecond(frequency)
	if ok {
		frequencySeconds := nano / int64(time.Second)
		return &frequencySeconds
	}
	return nil
}
func Label2Uns(labelDto *dao.UnsLabel) *dto.CreateTopicDto {
	unsDto := &dto.CreateTopicDto{}

	unsDto.ID = labelDto.ID
	unsDto.CreateAt = labelDto.CreateAt
	unsDto.UpdateAt = labelDto.UpdateAt
	unsDto.Flags = labelDto.WithFlags
	unsDto.Name = labelDto.LabelName
	unsDto.Path = "label/" + labelDto.LabelName
	unsDto.PathType = constants.PathTypeLabel

	// 状态设置：如果DelFlag为true则状态为0，否则为1
	//if labelDto.DelFlag != nil && *labelDto.DelFlag {
	//	unsDto.Status = 0
	//} else {
	//	unsDto.Status = 1
	//}

	freq := labelDto.SubscribeFrequency
	if freq != "" {
		unsDto.Frequency = freq
		unsDto.Protocol = map[string]interface{}{
			"frequency": freq,
		}
	}
	return unsDto
}
func Po2Dtos(poList []*dao.UnsNamespace) []*dto.CreateTopicDto {
	unsDtoList := make([]*dto.CreateTopicDto, len(poList))
	// BeanUtil.copyProperties equivalent - assuming a custom copy function
	copier.CopyWithOption(&unsDtoList, poList, copier.Option{IgnoreEmpty: true})
	for i, p := range poList {
		po2Dto(p, unsDtoList[i])
	}
	return unsDtoList
}
func Po2ApiDtos(poList []*dao.UnsNamespace) []*types.CreateTopicDto {
	if len(poList) == 0 {
		return nil
	}
	unsDtoList := make([]*types.CreateTopicDto, len(poList))
	// BeanUtil.copyProperties equivalent - assuming a custom copy function
	copier.CopyWithOption(&unsDtoList, poList, copier.Option{IgnoreEmpty: true})
	for i, p := range poList {
		Po2ApiDto(p, unsDtoList[i])
	}
	return unsDtoList
}
func Po2Dto(p *dao.UnsNamespace) *dto.CreateTopicDto {
	unsDto := &dto.CreateTopicDto{}
	copier.CopyWithOption(unsDto, p, copier.Option{IgnoreEmpty: true})
	po2Dto(p, unsDto)
	return unsDto
}
func po2Dto(p *dao.UnsNamespace, unsDto *dto.CreateTopicDto) {

	var withFlags int32
	if p.WithFlags != nil {
		withFlags = *p.WithFlags
	}
	unsDto.ID = p.ID
	unsDto.Flags = p.WithFlags
	unsDto.AddFlow = boPt(constants.WithFlow(withFlags))
	unsDto.AddDashBoard = boPt(constants.WithDashBoard(withFlags))
	unsDto.Save2DB = boPt(constants.WithSave2db(withFlags))
	unsDto.RetainTableWhenDeleteInstance = boPt(constants.WithRetainTableWhenDeleteInstance(withFlags))
	unsDto.ParentAlias = p.ParentAlias
	unsDto.ParentID = p.ParentID
	unsDto.Name = p.Name
	unsDto.LayRec = p.LayRec
	unsDto.ModelID = p.ModelID
	unsDto.ProtocolType = p.ProtocolType

	protocolStr := p.Protocol
	if protocolStr != "" && len(protocolStr) > 0 && protocolStr[0] == '{' {
		var protocol map[string]interface{}
		if err := JsonUtil.FromJson(protocolStr, &protocol); err == nil {
			if frequency, ok := protocol["frequency"].(string); ok {
				unsDto.FrequencySeconds = GetFrequencySeconds(frequency)
			}
			unsDto.Protocol = protocol
		}
	}

	jdbcType := common.GetSrcJdbcTypeByID(p.DataSrcID)
	unsDto.DataSrcID = jdbcType
	unsDto.Refers = p.Refers

	calculationExpr := p.Expression
	unsDto.Expression = calculationExpr
	//if calculationExpr != "" && compileExpression {
	//	dto.CompileExpression = ExpressionFunctions.CompileExpression(calculationExpr)
	//}

	fields := p.Fields
	unsDto.Fields = fields

	//if dto.DataType != nil && *dto.DataType == constants.AlarmRuleType && p.PathType == 2 {
	//	var ruleDefine AlarmRuleDefine
	//	if err := JsonUtil.FromJson(p.Protocol, &ruleDefine); err == nil {
	//		dto.AlarmRuleDefine = &ruleDefine
	//	}
	//}
	unsDto.ExtendFieldUsed = FieldUtils.ParseFlag(p.ExtendFieldFlags)
}
func Po2ApiDto(p *dao.UnsNamespace, unsDto *types.CreateTopicDto) {
	var withFlags int32
	if p.WithFlags != nil {
		withFlags = *p.WithFlags
	}
	unsDto.Id = p.ID
	unsDto.AddFlow = boPt(constants.WithFlow(withFlags))
	unsDto.AddDashBoard = boPt(constants.WithDashBoard(withFlags))
	unsDto.Save2db = boPt(constants.WithSave2db(withFlags))
	unsDto.RetainTableWhenDeleteInstance = boPt(constants.WithRetainTableWhenDeleteInstance(withFlags))
	unsDto.ParentAlias = p.ParentAlias
	unsDto.ParentId = p.ParentID
	unsDto.Name = p.Name
	unsDto.ModelId = p.ModelID

	calculationExpr := p.Expression
	unsDto.Expression = calculationExpr
	unsDto.ExtendFieldUsed = FieldUtils.ParseFlag(p.ExtendFieldFlags)
}
func boPt(b bool) *bool {
	return &b
}
func Dto2TreeResult(unsDto bo.NodeUnsInfo) *types.TopicTreeResult {
	result := &types.TopicTreeResult{}
	result.Id = strconv.FormatInt(unsDto.GetId(), 10)
	result.Alias = unsDto.GetAlias()
	if pid := unsDto.GetParentId(); pid != nil {
		strId := strconv.FormatInt(*pid, 10)
		result.ParentId = &strId
	}
	result.ParentAlias = unsDto.GetParentAlias()
	result.PathType = unsDto.GetPathType()
	name := PathUtil.GetName(unsDto.GetPath())
	result.Name = name
	result.Path = unsDto.GetPath()
	result.PathName = name
	result.DataType = unsDto.GetDataType()
	result.ParentDataType = unsDto.GetParentDataType()
	result.Mount = createMountDetailVo(unsDto)
	return result
}

func createMountDetailVo(unsDto bo.NodeUnsInfo) *types.MountDetailVo {
	if unsDto.GetMountType() == nil {
		return nil
	}

	mountDetailVo := &types.MountDetailVo{
		MountType:   unsDto.GetMountType(),
		MountSource: unsDto.GetMountSource(),
	}
	return mountDetailVo
}

var apiConvertOptions = copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
	{
		SrcType: copier.String,
		DstType: enums.FieldTypeInteger,
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := enums.GetFieldTypeByNameIgnoreCase(src.(string)); ok {
				return rs, nil
			}
			return nil, errors.Default
		},
	}, {
		SrcType: enums.FieldTypeInteger,
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(enums.FieldType); ok {
				return rs.Name(), nil
			}
			return nil, errors.Default
		},
	},
}}

func ConvertApiDto(apiDto types.CreateTopicDto) *dto.CreateTopicDto {
	var target dto.CreateTopicDto
	copier.CopyWithOption(&target, apiDto, apiConvertOptions)
	return &target
}
func ConvertApiUpdateDto(apiDto *types.UpdateUnsDto) *dto.CreateTopicDto {
	var target dto.CreateTopicDto
	copier.CopyWithOption(&target, apiDto, apiConvertOptions)
	return &target
}

func ConvertApiDtos(apiDto []types.CreateTopicDto) (target []*dto.CreateTopicDto) {
	copier.CopyWithOption(&target, apiDto, apiConvertOptions)
	return target
}
func ConvertFields(apiDto []*dto.FieldDefine) (target []*types.FieldDefine) {
	copier.CopyWithOption(&target, apiDto, apiConvertOptions)
	return target
}
