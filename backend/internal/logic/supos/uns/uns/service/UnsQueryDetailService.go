package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strconv"
	"strings"
)

func (l *UnsQueryService) GetInstanceDetail(ctx context.Context, req *types.InstanceDetailReq, alias string) (resp *types.InstanceDetailResp, err error) {
	detail := &types.InstanceDetail{}
	db := dao.GetDb(ctx)
	var po *dao.UnsNamespace
	if id := req.Id; id > 0 {
		po, err = l.unsMapper.SelectById(db, id)
	} else if len(alias) > 0 {
		po, err = l.unsMapper.GetByAlias(db, alias)
	}
	resp = &types.InstanceDetailResp{}
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return
	} else if po == nil {
		resp.Code = 200
		resp.Msg = I18nUtils.GetMessage("uns.file.not.found")
		return
	}
	detail.SubscribeEnable = constants.WithSubscribeEnable(po.WithFlags)
	protocol := po.Protocol
	if detail.SubscribeEnable && strings.HasPrefix(protocol, "{") {
		var pMap = make(map[string]string)
		JsonUtil.FromJson(protocol, &pMap)
		//if freq, has := pMap["frequency"]; has {
		//	detail.SubscribeFrequency = freq
		//}
	}
	detail.Id = strconv.FormatInt(po.Id, 10)
	UnsConverter.CopyProperties(po, detail)
	detail.Topic = base.SanYuan(constants.UseAliasAsTopic, po.Alias, po.Path)
	if fs := detail.Fields; len(fs) > 0 {
		for _, f := range fs {
			f.Index = nil
		}
	}
	if templateId := po.ModelId; templateId != nil {
		if template, er := l.unsMapper.SelectById(db, *templateId); er == nil && template != nil {
			detail.ModelId = strconv.FormatInt(*templateId, 10)
			detail.ModelId = template.Name
			detail.TemplateAlias = template.Alias
		}
	}
	if ms := l.mountService; ms != nil {
		detail.Mount = ms.ParseMountDetail(po, false)
	}
	resp.Data = detail
	return
}
func (l *UnsQueryService) GetModelDefinition(ctx context.Context, req *types.ModelDetailReq, alias string) (resp *types.ModelDetailResp, err error) {
	dto := &types.ModelDetail{}
	db := dao.GetDb(ctx)
	var po *dao.UnsNamespace
	if id := req.Id; id > 0 {
		po, err = l.unsMapper.SelectById(db, id)
	} else if len(alias) > 0 {
		po, err = l.unsMapper.GetByAlias(db, alias)
	}
	resp = &types.ModelDetailResp{}
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return
	} else if po == nil {
		resp.Code = 200
		resp.Msg = I18nUtils.GetMessage("uns.model.not.found")
		return
	}
	dto.SubscribeEnable = constants.WithSubscribeEnable(po.WithFlags)
	protocol := po.Protocol
	if dto.SubscribeEnable && strings.HasPrefix(protocol, "{") {
		var pMap = make(map[string]string)
		JsonUtil.FromJson(protocol, &pMap)
		if freq, has := pMap["frequency"]; has {
			dto.SubscribeFrequency = freq
		}
	}
	dto.Id = strconv.FormatInt(po.Id, 10)
	UnsConverter.CopyProperties(po, dto)
	dto.Topic = base.SanYuan(constants.UseAliasAsTopic, po.Alias, po.Path)
	if fs := dto.Fields; len(fs) > 0 {
		for _, f := range fs {
			f.Index = nil
		}
	}
	if templateId := po.ModelId; templateId != nil {
		if template, er := l.unsMapper.SelectById(db, *templateId); er == nil && template != nil {
			dto.ModelId = strconv.FormatInt(*templateId, 10)
			dto.ModelName = template.Name
			dto.TemplateAlias = template.Alias
		}
	}
	if ms := l.mountService; ms != nil {
		dto.Mount = ms.ParseMountDetail(po, false)
	}
	resp.Data = dto
	return
}

/*func (l *UnsQueryService) setDetailInfo(file bo.UnsInfo, dto *types.InstanceDetail, setMount bool) {
	fs := file.Refers
	var origPo *CreateTopicDto

	if len(fs) > 0 {
		if file.DataType != nil && *file.DataType == Constants.CITING_TYPE {
			if len(fs) > 0 && fs[0].ID != nil {
				orig := l.unsDefinitionService.GetDefinitionById(*fs[0].ID)
				if orig != nil {
					origPo = orig
				}
			}
		}
	}

	// 确定目标UNS
	unsTarget := file
	if origPo != nil {
		unsTarget = origPo
	}

	// 设置字段
	if fields := unsTarget.GetFields(); fields != nil {
		fieldDefines := l.getDisplayFields(unsTarget, unsTarget.Fields)
		dto.Fields = fieldDefines
	}

	// 设置基本信息
	if id := file.GetId(); id > 0 {
		dto.Id = strconv.FormatInt(id, 10)
	}
	dto.DataType = file.GetDataType()
	dto.Alias = file.GetAlias()
	dto.Path = file.GetPath()

	if constants.UseAliasAsTopic {
		dto.Topic = file.Alias
	} else {
		dto.Topic = file.Path
	}

	dto.DataPath = file.DataPath
	dto.Protocol = file.ProtocolMap

	// 设置引用和表达式
	if len(fs) > 0 {
		l.setRefersAndExpression(fs, unsTarget.Expression, unsTarget.CalculationType, dto.Protocol, dto)
	}

	dto.CalculationType = unsTarget.CalculationType
	dto.Description = file.Description
	dto.CreateTime = l.getDatetime(file.CreateAt)
	dto.UpdateTime = l.getDatetime(file.UpdateAt)
	dto.Alias = file.Alias
	dto.Name = file.Name
	dto.DisplayName = file.DisplayName
	dto.PathName = PathUtil.GetName(file.Path)
	dto.Extend = file.Extend

	// 设置读写模式、保存到数据库、扩展字段使用、挂载信息
	dto.ReadWriteMode = unsTarget.ReadWriteMode
	dto.ExtendFieldUsed = FieldUtils.ParseFlag(unsTarget.ExtendFieldFlags)

	if setMount {
		dto.Mount = l.mountService.ParseMountDetail(unsTarget, false)
	}

	// 设置标志位
	if unsTarget.Flags != nil {
		flags := *unsTarget.Flags
		dto.WithFlow = Constants.WithFlow(flags)
		dto.WithDashboard = Constants.WithDashBoard(flags)
		dto.WithSave2db = Constants.WithSave2db(flags)
		dto.Save2db = Constants.WithSave2db(flags)
		dto.SubscribeEnable = Constants.WithSubscribeEnable(flags)
	}

	// 设置标签列表
	if file.LabelIds != nil && len(file.LabelIds) > 0 {
		labelIds := make([]int64, 0, len(file.LabelIds))
		for id := range file.LabelIds {
			labelIds = append(labelIds, id)
		}

		labelPos, err := l.labelMapper.SelectByIds(labelIds)
		if err == nil {
			labelList := make([]*UnsLabelVo, 0, len(labelPos))
			for _, p := range labelPos {
				vo := &UnsLabelVo{}
				// BeanUtils.copyProperties equivalent
				l.copyProperties(p, vo)
				if p.ID != nil {
					vo.ID = *p.ID
				}
				labelList = append(labelList, vo)
			}
			dto.LabelList = labelList
		}
	}

	// 设置模板信息
	if file.ModelId != nil {
		template := l.unsDefinitionService.GetDefinitionById(*file.ModelId)
		if template != nil {
			dto.ModelId = strconv.FormatInt(*file.ModelId, 10)
			dto.ModelName = template.Name
			dto.TemplateName = template.Name
			dto.TemplateAlias = template.Alias
		}
	}
}
*/
