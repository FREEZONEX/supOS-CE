package UnsConverter

import (
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/utils/FieldFlags"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"fmt"
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
func Label2Uns(labelDto *dao.UnsLabel) *types.CreateTopicDto {
	unsDto := &types.CreateTopicDto{}

	unsDto.Id = labelDto.ID
	unsDto.CreateAt = labelDto.CreateAt.UnixMilli()
	unsDto.UpdateAt = labelDto.UpdateAt.UnixMilli()
	unsDto.WithFlags = labelDto.WithFlags
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
func Po2Dtos(poList []*dao.UnsNamespace) []*types.CreateTopicDto {
	unsDtoList := make([]*types.CreateTopicDto, len(poList))
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
func Po2Dto(p *dao.UnsNamespace) *types.CreateTopicDto {
	unsDto := &types.CreateTopicDto{}
	copier.CopyWithOption(unsDto, p, copier.Option{IgnoreEmpty: true})
	po2Dto(p, unsDto)
	return unsDto
}
func po2Dto(p *dao.UnsNamespace, unsDto *types.CreateTopicDto) {

	var withFlags int32
	if p.WithFlags != nil {
		withFlags = *p.WithFlags
	}
	unsDto.Id = p.Id
	unsDto.Alias = p.Alias
	unsDto.DataType = p.DataType
	unsDto.ParentDataType = p.ParentDataType
	if tbn := p.TableName_; tbn != nil && *tbn != "" {
		unsDto.TableName = *tbn
	} else {
		unsDto.TableName = ""
	}
	unsDto.WithFlags = p.WithFlags
	unsDto.AddFlow = boPt(constants.WithFlow(withFlags))
	unsDto.AddDashBoard = boPt(constants.WithDashBoard(withFlags))
	unsDto.Save2Db = boPt(constants.WithSave2db(withFlags))
	unsDto.RetainTableWhenDeleteInstance = boPt(constants.WithRetainTableWhenDeleteInstance(withFlags))
	unsDto.AccessLevel = constants.WithReadOnly(withFlags)
	unsDto.ParentAlias = p.ParentAlias
	unsDto.ParentId = p.ParentId
	unsDto.Name = p.Name
	unsDto.LayRec = p.LayRec
	unsDto.ModelId = p.ModelId
	unsDto.ProtocolType = p.ProtocolType

	protocolStr := ""
	if p.Protocol != nil {
		protocolStr = *p.Protocol
	}
	if len(protocolStr) > 2 && protocolStr[0] == '{' {
		var protocol map[string]interface{}
		if err := JsonUtil.FromJson(protocolStr, &protocol); err == nil && len(protocol) > 0 {
			if frequency, ok := protocol["frequency"].(string); ok {
				unsDto.FrequencySeconds = GetFrequencySeconds(frequency)
			}
			if jsf, has := protocol["jsf"]; has {
				if str, isStr := jsf.(string); isStr {
					JsonUtil.FromJson(str, &unsDto.JsonFields)
				} else {
					JsonUtil.FromJson(fmt.Sprint(jsf), &unsDto.JsonFields)
				}
			}
			unsDto.Protocol = protocol
		}
	}

	unsDto.DataSrcID = p.DataSrcId
	unsDto.Refers = p.Refers

	calculationExpr := p.Expression
	unsDto.Expression = calculationExpr
	//if calculationExpr != "" && compileExpression {
	//	dto.CompileExpression = ExpressionFunctions.CompileExpression(calculationExpr)
	//}

	fields := p.Fields
	unsDto.Fields = fields
	if len(fields) > 0 {
		for _, field := range fields {
			if field.IsSystemField() {
				field.SystemField = base.OptionalTrue
			}
		}
	}
	if labels := p.LabelIds; len(labels) > 0 {
		unsDto.LabelNames = base.MapValues(labels)
	}
	//if dto.DataType != nil && *dto.DataType == constants.AlarmRuleType && p.PathType == 2 {
	//	var ruleDefine AlarmRuleDefine
	//	if err := JsonUtil.FromJson(p.Protocol, &ruleDefine); err == nil {
	//		dto.AlarmRuleDefine = &ruleDefine
	//	}
	//}
	unsDto.ExtendFieldUsed = FieldFlags.ParseFlag(p.ExtendFieldFlags)
}
func Po2ApiDto(p *dao.UnsNamespace, unsDto *types.CreateTopicDto) {
	var withFlags int32
	if p.WithFlags != nil {
		withFlags = *p.WithFlags
	}
	unsDto.Id = p.Id
	unsDto.AddFlow = boPt(constants.WithFlow(withFlags))
	unsDto.AddDashBoard = boPt(constants.WithDashBoard(withFlags))
	unsDto.Save2Db = boPt(constants.WithSave2db(withFlags))
	unsDto.RetainTableWhenDeleteInstance = boPt(constants.WithRetainTableWhenDeleteInstance(withFlags))
	unsDto.ParentAlias = p.ParentAlias
	unsDto.ParentId = p.ParentId
	unsDto.Name = p.Name
	unsDto.ModelId = p.ModelId

	calculationExpr := p.Expression
	unsDto.Expression = calculationExpr
	unsDto.ExtendFieldUsed = FieldFlags.ParseFlag(p.ExtendFieldFlags)
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
	mt := unsDto.GetMountType()
	if mt == nil || *mt == 0 {
		return nil
	}

	mountDetailVo := &types.MountDetailVo{
		MountType:   mt,
		MountSource: unsDto.GetMountSource(),
	}
	return mountDetailVo
}

var int64Temp = int64(0)
var apiConvertOptions = copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
	{
		SrcType: copier.String,
		DstType: types.FieldTypeInteger,
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := types.GetFieldTypeByNameIgnoreCase(src.(string)); ok {
				return rs, nil
			}
			return nil, errors.Default
		},
	}, {
		SrcType: types.FieldTypeInteger,
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(types.FieldType); ok {
				return rs.Name(), nil
			}
			return nil, errors.Default
		},
	}, {
		SrcType: time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(time.Time); ok {
				return rs.UnixMilli(), nil
			}
			return nil, errors.Default
		},
	}, {
		SrcType: &time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(*time.Time); ok && rs != nil {
				return rs.UnixMilli(), nil
			} else if ok && rs == nil {
				return int64(0), nil
			}
			return nil, errors.Default
		},
	}, {
		SrcType: int64(0),
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(int64); ok {
				return strconv.FormatInt(rs, 10), nil
			}
			return nil, errors.Default
		},
	}, {
		SrcType: copier.String,
		DstType: int64(0),
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(string); ok {
				return strconv.ParseInt(rs, 10, 64)
			}
			return nil, errors.Default
		},
	}, {
		SrcType: copier.String,
		DstType: &int64Temp,
		Fn: func(src interface{}) (dst interface{}, err error) {
			if rs, ok := src.(string); ok {
				var num int64
				num, err = strconv.ParseInt(rs, 10, 64)
				dst = num
			}
			return nil, errors.Default
		},
	},
}}

func ConvertApiUpdateDto(apiDto *types.UpdateUnsDto) *types.CreateTopicDto {
	var target types.CreateTopicDto
	er := copier.CopyWithOption(&target, apiDto, copier.Option{IgnoreEmpty: true})
	if er != nil {
		target.Id = apiDto.Id
		target.Alias = apiDto.Alias
		target.ParentAlias = apiDto.ParentAlias
		target.ParentId = apiDto.ParentId
		target.Fields = apiDto.Fields
		target.Name = apiDto.Name
		target.ModelId = apiDto.ModelId
		target.ModelAlias = apiDto.ModelAlias
		target.Description = apiDto.Description
		target.WithFlags = apiDto.WithFlags
		target.AddFlow = apiDto.AddFlow
		target.AddDashBoard = apiDto.AddDashBoard
		target.SubscribeEnable = apiDto.SubscribeEnable
	}
	return &target
}

func CopyProperties(from any, to any) error {
	return copier.CopyWithOption(to, from, apiConvertOptions)
}
func CopyPropertiesDefault(from any, to any) error {
	return copier.Copy(to, from)
}
func LabelPo2Vo(po *dao.UnsLabel) (vo *types.LabelVo) {
	vo = &types.LabelVo{ID: po.ID, LabelName: po.LabelName, CreateTime: po.CreateAt.UnixMilli()}
	if po.SubscribeAt != nil && !po.SubscribeAt.IsZero() {
		vo.SubscribeAt = po.SubscribeAt.UnixMilli()
	}
	vo.Topic = "label/" + po.LabelName
	if flags := base.P2v(po.WithFlags); flags > 0 {
		enable := constants.WithSubscribeEnable(flags)
		vo.SubscribeEnable = &enable
		if enable {
			vo.SubscribeFrequency = po.SubscribeFrequency
		}
	}
	return vo
}
