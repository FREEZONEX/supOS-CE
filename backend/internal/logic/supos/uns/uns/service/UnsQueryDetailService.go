package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/FieldUtils"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	"backend/internal/logic/supos/uns/uns/bo"
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
	l.setDetailInfo(ctx, po, detail, true)
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
	dto.SubscribeEnable = constants.WithSubscribeEnable(base.P2v(po.WithFlags))
	protocol := po.Protocol
	if dto.SubscribeEnable && protocol != nil && strings.HasPrefix(*protocol, "{") {
		var pMap = make(map[string]string)
		JsonUtil.FromJson(*protocol, &pMap)
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
	if ms := l.mountService; ms != nil {
		dto.Mount = ms.ParseMountDetail(po, false)
	}
	resp.Data = dto
	return
}

func (l *UnsQueryService) setDetailInfo(ctx context.Context, file bo.UnsInfo, dto *types.InstanceDetail, setMount bool) {
	fs := file.GetRefers()
	var origPo *types.CreateTopicDto
	db := dao.GetDb(ctx)
	ctx = dao.SetDb(ctx, db)
	if len(fs) > 0 {
		if dataType := file.GetDataType(); dataType != nil && *dataType == constants.CitingType {
			if len(fs) > 0 && fs[0].Id > 0 {
				orig, _ := l.unsMapper.SelectById(db, fs[0].Id)
				if orig != nil {
					origPo = UnsConverter.Po2Dto(orig)
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
		fieldDefines := getDisplayFields(unsTarget, unsTarget.GetFields())
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
		dto.Topic = file.GetAlias()
	} else {
		dto.Topic = file.GetPath()
	}
	if dp := file.GetDataPath(); len(dp) > 0 {
		dto.DataPath = &dp
	}
	dto.Protocol = file.GetProtocolMap()

	// 设置引用和表达式
	if len(fs) > 0 {
		l.calcService.setRefersAndExpression(fs, unsTarget.GetExpression(), unsTarget.GetCalculationType(), dto.Protocol, dto)
	}

	//dto.CalculationType = unsTarget.CalculationType
	dto.Description = file.GetDescription()
	dto.CreateTime = file.GetCreateAt()
	dto.UpdateTime = file.GetUpdateAt()
	dto.Alias = file.GetAlias()
	dto.Name = file.GetName()
	dto.DisplayName = file.GetDisplayName()
	dto.PathName = PathUtil.GetName(file.GetPath())
	dto.Extend = file.GetExtend()

	// 设置读写模式、保存到数据库、扩展字段使用、挂载信息
	//dto.ReadWriteMode = unsTarget.ReadWriteMode
	//dto.ExtendFieldUsed = FieldFlags.ParseFlag(unsTarget.GetExtendFieldFlags())

	if mc := l.mountService; setMount && mc != nil {
		dto.Mount = mc.ParseMountDetail(unsTarget, false)
	}

	// 设置标志位
	if flagsP := unsTarget.GetFlags(); flagsP != nil {
		flags := *flagsP
		dto.WithFlow = constants.WithFlow(flags)
		dto.WithDashboard = constants.WithDashBoard(flags)
		dto.WithSave2db = constants.WithSave2db(flags)
		dto.Save2db = constants.WithSave2db(flags)
		dto.SubscribeEnable = constants.WithSubscribeEnable(flags)
	}

	// 设置标签列表
	if labelIds := file.GetLabelIds(); len(labelIds) > 0 {
		labelPos, _ := l.labelMapper.ListByIds(db, base.MapKeys(labelIds))
		if len(labelPos) > 0 {
			dto.LabelList = base.Map[*dao.UnsLabel, types.LabelVo](labelPos, func(e *dao.UnsLabel) (rs types.LabelVo) {
				rs = *UnsConverter.LabelPo2Vo(e)
				return rs
			})
		}
	}

	// 设置模板信息
	if templateId := file.GetModelId(); templateId != nil {
		if template, er := l.unsMapper.SelectById(db, *templateId); er == nil && template != nil {
			dto.ModelId = strconv.FormatInt(*templateId, 10)
			dto.ModelId = template.Name
			dto.TemplateName = template.Name
			dto.TemplateAlias = template.Alias
		}
	}
}
func getDisplayFields(unsInfo bo.UnsInfo, fields []*types.FieldDefine) []*types.FieldDefine {
	dataType := unsInfo.GetDataType()
	if dataType == nil {
		return fields
	}
	if *dataType == constants.TimeSequenceType || *dataType == constants.CalculationRealType {
		return filterFieldsForTimeSequence(fields)
	} else {
		return filterFieldsForOtherTypes(unsInfo, fields)
	}
}

func filterFieldsForTimeSequence(fields []*types.FieldDefine) []*types.FieldDefine {
	result := make([]*types.FieldDefine, 0, len(fields))

	for _, fd := range fields {
		name := fd.GetName()
		tbValueName := fd.GetTbValueName()

		// 保留不包含系统字段前缀且没有表值名称的字段
		if !strings.HasPrefix(name, constants.SystemFieldPrev) && tbValueName == nil {
			result = append(result, fd)
		}
	}

	return result
}

func filterFieldsForOtherTypes(unsInfo bo.UnsInfo, fields []*types.FieldDefine) []*types.FieldDefine {
	jdbcType := unsInfo.GetSrcJdbcType()
	if jdbcType == 0 {
		return fields
	}

	ct := FieldUtils.GetTimestampField(fields)
	qos := FieldUtils.GetQualityField(fields, jdbcType.TypeCode())

	result := make([]*types.FieldDefine, 0, len(fields))

	for _, fd := range fields {
		// 跳过时间戳字段、质量字段和系统字段
		if fd == ct || fd == qos || strings.HasPrefix(fd.GetName(), constants.SystemFieldPrev) {
			continue
		}
		result = append(result, fd)
	}

	return result
}
