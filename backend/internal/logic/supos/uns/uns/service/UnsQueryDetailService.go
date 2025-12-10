package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/FieldUtils"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strconv"
	"strings"

	"gorm.io/gorm"
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
	dto := &types.ModelDetail{}
	l.setDetailInfo(ctx, po, dto, false)
	resp.Code, resp.Data = 200, dto
	return
}

func (l *UnsQueryService) setDetailInfo(ctx context.Context, file types.UnsInfo, dto bo.UnsDetail, setMount bool) {
	fs := file.GetRefers()
	var origPo *types.CreateTopicDto
	var db *gorm.DB
	if len(fs) > 0 {
		db = dao.GetDb(ctx)
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
	dto.SetPathType(unsTarget.GetPathType())
	// 设置字段
	if fields := unsTarget.GetFields(); len(fields) > 0 {
		fieldDefines := getDisplayFields(unsTarget, fields)
		dto.SetFields(fieldDefines)
	}

	// 设置基本信息
	if id := file.GetId(); id > 0 {
		dto.SetId(strconv.FormatInt(id, 10))
	}
	dto.SetDataType(file.GetDataType())
	dto.SetParentDataType(file.GetParentDataType())
	dto.SetAlias(file.GetAlias())
	dto.SetPath(file.GetPath())
	if constants.UseAliasAsTopic {
		dto.SetTopic(file.GetAlias())
	} else {
		dto.SetTopic(file.GetPath())
	}
	if dp := file.GetDataPath(); len(dp) > 0 {
		dto.SetDataPath(&dp)
	}
	// 设置引用和表达式
	if len(fs) > 0 {
		l.calcService.setRefersAndExpression(fs, unsTarget.GetExpression(), unsTarget.GetCalculationType(), file.GetProtocolMap(), dto)
	}

	//dto.CalculationType = unsTarget.CalculationType
	dto.SetDescription(file.GetDescription())
	dto.SetCreateTime(file.GetCreateAt())
	dto.SetUpdateTime(file.GetUpdateAt())
	dto.SetName(file.GetName())
	dto.SetDisplayName(file.GetDisplayName())
	dto.SetPathName(PathUtil.GetName(file.GetPath()))
	dto.SetExtend(file.GetExtend())

	// 设置读写模式、保存到数据库、扩展字段使用、挂载信息
	//dto.ReadWriteMode = unsTarget.ReadWriteMode
	//dto.ExtendFieldUsed = FieldFlags.ParseFlag(unsTarget.GetExtendFieldFlags())

	if mc := l.mountService; setMount && mc != nil {
		dto.SetMount(mc.ParseMountDetail(unsTarget, false))
	}
	dto.SetTable(file.GetTable())
	dto.SetTbFieldName(file.GetTbFieldName())
	// 设置标志位
	if flagsP := unsTarget.GetFlags(); flagsP != nil {
		flags := *flagsP
		dto.SetWithFlow(constants.WithFlow(flags))
		dto.SetWithDashboard(constants.WithDashBoard(flags))
		dto.SetWithSave2db(constants.WithSave2db(flags))
		dto.SetSave2db(constants.WithSave2db(flags))
		dto.SetSubscribeEnable(constants.WithSubscribeEnable(flags))
		dto.SetAccessLevel(constants.WithReadOnly(flags))
	}
	dto.SetProtocol(file.GetProtocolMap())
	// 设置标签列表
	if labelIds := file.GetLabelIds(); len(labelIds) > 0 {
		if db == nil {
			db = dao.GetDb(ctx)
		}
		labelPos, _ := l.labelMapper.ListByIds(db, base.MapKeys(labelIds))
		if len(labelPos) > 0 {
			dto.SetLabelList(base.Map[*dao.UnsLabel, types.LabelVo](labelPos, func(e *dao.UnsLabel) (rs types.LabelVo) {
				rs = *UnsConverter.LabelPo2Vo(e)
				return rs
			}))
		}
	}

	// 设置模板信息
	if templateId := file.GetModelId(); templateId != nil {
		if db == nil {
			db = dao.GetDb(ctx)
		}
		if template, er := l.unsMapper.SelectById(db, *templateId); er == nil && template != nil {
			dto.SetModelId(strconv.FormatInt(*templateId, 10))
			dto.SetModelName(template.Name)
			dto.SetTemplateName(template.Name)
			dto.SetTemplateAlias(template.Alias)
		}
	}
}
func getDisplayFields(unsInfo types.UnsInfo, fields []*types.FieldDefine) []*types.FieldDefine {
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
			if fd.IsSystemField() {
				fd.SystemField = base.OptionalTrue
			}
		}
	}

	return result
}

func filterFieldsForOtherTypes(unsInfo types.UnsInfo, fields []*types.FieldDefine) []*types.FieldDefine {
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
		if fd.IsSystemField() {
			fd.SystemField = base.OptionalTrue
		}
		result = append(result, fd)
	}

	return result
}
