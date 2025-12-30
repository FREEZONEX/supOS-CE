package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/LeastTopNodeUtil"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/utils/FieldFlags"
	"backend/internal/common/utils/FieldUtils"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

func checkInstanceFields(modelFields []*types.FieldDefine, insFields []*types.FieldDefine) string {
	if modelFields == nil {
		modelFields = make([]*types.FieldDefine, 0)
	}
	if insFields == nil {
		insFields = make([]*types.FieldDefine, 0)
	}

	insMap := make(map[string]*types.FieldDefine, len(insFields))
	for _, insField := range insFields {
		name := insField.Name
		if !insField.IsSystemField() {
			if _, exists := insMap[name]; exists {
				return "fields name duplicate: " + name
			}
			insMap[name] = insField
		}
	}

	for _, mf := range modelFields {
		name := mf.Name
		if !mf.IsSystemField() {
			insF, exists := insMap[name]
			if !exists {
				return "instance need field: " + name
			} else if mf.Type != insF.Type {
				return fmt.Sprintf("instance field type changed: %s, %s -> %s", name, mf.Type, insF.Type)
			}
			delete(insMap, name)
		}
	}

	if len(insMap) > 0 {
		var unknownFields []string
		for name := range insMap {
			unknownFields = append(unknownFields, name)
		}
		return "instance has unknown Fields in model: " + strings.Join(unknownFields, ", ")
	}
	return ""
}

func initParamsUns(topicDtos []*types.CreateTopicDto, errTipMap map[string]string) (pathMap map[int16]map[string]*types.CreateTopicDto) {
	pathMap = make(map[int16]map[string]*types.CreateTopicDto, 4)
	for _, topicDto := range topicDtos {
		checkTopicDto(errTipMap, pathMap, func(dto *types.CreateTopicDto) {
			vMap, has := pathMap[dto.PathType]
			if !has {
				var initSize = 16
				switch dto.PathType {
				case constants.PathTypeFile:
					initSize = len(topicDtos)
				case constants.PathTypeDir:
					initSize = 32
				case constants.PathTypeTemplate:
					initSize = 65
				}
				vMap = make(map[string]*types.CreateTopicDto, initSize)
				pathMap[dto.PathType] = vMap
			}
			vMap[dto.Alias] = dto
		}, topicDto)
	}
	return pathMap
}
func addAlias(bos []*types.CreateTopicDto, aliasSet map[string]bool, ids map[int64]bool) {
	if len(bos) == 0 {
		return
	}
	for _, unsDto := range bos {
		// 添加alias
		if alias := unsDto.Alias; alias != "" {
			aliasSet[alias] = true
		}

		// 添加referUns
		if refAlias := unsDto.ReferUns; refAlias != "" {
			aliasSet[refAlias] = true
		}

		// 添加modelAlias
		if modelAlias := unsDto.ModelAlias; modelAlias != nil {
			aliasSet[*modelAlias] = true
		}

		// 添加parentAlias
		if folderAlias := unsDto.ParentAlias; folderAlias != nil {
			aliasSet[*folderAlias] = true
		}

		// 处理referIds
		if referIds := unsDto.ReferIds; len(referIds) > 0 {
			// 添加所有referIds到ids集合
			for _, id := range referIds {
				ids[id] = true
			}

			// 如果refers为空，根据referIds创建InstanceField数组
			if len(unsDto.Refers) == 0 {
				refers := make([]*types.InstanceField, len(referIds))
				for i, id := range referIds {
					refers[i] = &types.InstanceField{Id: id, Alias: ""}
				}
				unsDto.Refers = refers
			}
		}

		// 添加各种ID到ids集合
		if unsId := unsDto.Id; unsId != 0 {
			ids[unsId] = true
		}
		if pid := unsDto.ParentId; pid != nil {
			ids[*pid] = true
		}
		if mid := unsDto.ModelId; mid != nil {
			ids[*mid] = true
		}

		// 处理refers中的字段
		if refers := unsDto.Refers; len(refers) > 0 {
			for _, field := range refers {
				if id := field.Id; id != 0 {
					ids[id] = true
				}
				if alias := field.Alias; alias != "" {
					aliasSet[alias] = true
				}
			}
		}
	}
}

func eqStrP(s1, s2 *string) bool {
	if s1 == nil && s2 == nil {
		return true
	} else if s1 == nil || s2 == nil {
		return false
	}
	return *s1 == *s2
}
func tryFillIdOrAlias(paramFiles map[string]*types.CreateTopicDto, existsUns map[string]*dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace, errTipMap map[string]string) {
	for key, topicDto := range paramFiles {
		// 处理主对象的ID和Alias
		id := topicDto.Id
		alias := topicDto.Alias
		if id != 0 && alias == "" {
			if po, exists := dbFiles[id]; exists {
				topicDto.Alias = po.Alias
			}
		} else if id == 0 && alias != "" {
			if po, exists := existsUns[alias]; exists {
				topicDto.Id = po.Id
			}
		}

		// 处理父级ID和ParentAlias
		pid := topicDto.ParentId
		parentAlias := topicDto.ParentAlias
		if pid == nil && parentAlias != nil {
			if parent, exists := existsUns[*parentAlias]; exists {
				topicDto.ParentId = &parent.Id
			}
		} else if parentAlias == nil && pid != nil {
			if parent, exists := dbFiles[*pid]; exists {
				topicDto.ParentAlias = &parent.Alias
			}
		}

		// 处理引用字段
		refers := topicDto.Refers
		if len(refers) > 0 {
			for _, field := range refers {
				refID := field.Id
				refAlias := field.Alias
				if refID == 0 && refAlias != "" {
					if refPo, exists := existsUns[refAlias]; exists {
						field.Id = refPo.Id
					} else {
						// 删除当前元素并记录错误
						delete(paramFiles, key)
						errTipMap[topicDto.GainBatchIndex()] = I18nUtils.GetMessage("uns.topic.calc.expression.topic.ref.notFound", topicDto.Alias)
						break // 跳出内层循环，继续处理下一个元素
					}
				} else if refID != 0 && refAlias == "" {
					if refPo, exists := dbFiles[refID]; exists {
						field.Alias = refPo.Alias
					} else {
						delete(paramFiles, key)
						errTipMap[topicDto.GainBatchIndex()] = I18nUtils.GetMessage("uns.topic.calc.expression.topic.ref.notFound", topicDto.Alias)
						break
					}
				}
			}
		}
	}
}

// 添加数据库PO到映射
func addDbPo(unsPos []*dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace, aliasMap map[string]*dao.UnsNamespace) {
	for _, po := range unsPos {
		putTemp(dbFiles, aliasMap, po)
	}
}

var _ZERO = int16(0)
var LOGIC_REMOVED = int16(0)
var OK = int16(1)

// 临时处理PO
func putTemp(dbFiles map[int64]*dao.UnsNamespace, aliasMap map[string]*dao.UnsNamespace, po *dao.UnsNamespace) {
	if base.P2v(po.Status) == LOGIC_REMOVED {
		// 创建新的PO对象，只保留ID和Alias
		newPo := &dao.UnsNamespace{
			Id:     po.Id,
			Alias:  po.Alias,
			Status: &LOGIC_REMOVED,
		}
		// 使用新对象替换原对象
		po = newPo
	}

	// 添加到映射
	aliasMap[po.Alias] = po
	dbFiles[po.Id] = po
}

func checkTopicDto(errTipMap map[string]string,
	pathMap map[int16]map[string]*types.CreateTopicDto,
	put func(*types.CreateTopicDto),
	d *types.CreateTopicDto) {

	pathType := d.PathType
	//TODO 假设 validator 在 Go 中已实现
	//violations := validator.Validate(dto)
	batchIndex := d.GainBatchIndex()

	//if len(violations) > 0 {
	//	er := strings.Builder{}
	//	er.Grow(128)
	//	addValidErrMsg(&er, violations)
	//	errTipMap[batchIndex] = er.String()
	//	return
	//}

	alias := d.Alias
	for _, mp := range pathMap {
		if v, has := mp[alias]; has {
			dupUns, _ := json.Marshal(v)
			curUns, _ := json.Marshal(d)
			msg := I18nUtils.GetMessage("uns.alias.duplicate") + ": " + string(dupUns) + ", alias=" + alias + ", curUns=" + string(curUns)
			errTipMap[batchIndex] = msg
			return
		}
	}
	if alias != "" && alias == base.P2v(d.ParentAlias) {
		errTipMap[batchIndex] = I18nUtils.GetMessage("uns.circularDependency")
		return
	}
	if pathType == constants.PathTypeDir { // 当前是文件夹
		put(d)
	} else if pathType == constants.PathTypeFile { // 当前是文件
		dataType := d.DataType
		if dataType == nil {
			if d.Id == 0 {
				msg := I18nUtils.GetMessage("uns.file.dataType.empty")
				errTipMap[batchIndex] = msg
				return
			}
		} else if !constants.IsValidDataType(*dataType) {
			msg := fmt.Sprint(I18nUtils.GetMessage("uns.file.dataType.invalid"), *dataType)
			errTipMap[batchIndex] = msg
			return
		}

		fields := d.Fields
		if len(fields) == 0 && base.P2v(dataType) == constants.MergeType {
			mergeMaxLen := 512 * 1024
			mergeField := &types.FieldDefine{
				Name:   "data_json",
				Type:   types.FieldTypeString,
				MaxLen: &mergeMaxLen, // 聚合的字段总长度限制改大，不能超过mqtt消息长度限制
			}
			fields = []*types.FieldDefine{mergeField}
			d.Fields = fields
		}
		put(d)
	} else if pathType == constants.PathTypeTemplate { // 当前是模版
		d.DataType = &_ZERO
		d.ParentId = &templateRootID
		d.ParentAlias = &templateRootAlias
		put(d)
	} else {
		errTipMap[batchIndex] = I18nUtils.GetMessage("uns.import.type.error")
		return
	}

	if d.Frequency != "" {
		protocol := d.Protocol
		if protocol == nil {
			protocol = make(map[string]interface{})
		}
		frequency := d.Frequency
		protocol["frequency"] = frequency
		d.Protocol = protocol
		d.FrequencySeconds = UnsConverter.GetFrequencySeconds(frequency)
	}
}

var templateRootID = int64(0)
var templateRootAlias = "__templates__"

func setJdbcType(unsDto *types.CreateTopicDto) {
	dataType := unsDto.DataType
	jdbcType := unsDto.DataSrcID
	if jdbcType == 0 && dataType != nil && unsDto.PathType == constants.PathTypeFile {
		switch *dataType {
		case constants.CalculationHistType, constants.CalculationRealType, constants.TimeSequenceType:
			jdbcType = types.SrcJdbcTypeTimeScaleDB.Id()
		case constants.AlarmRuleType, constants.RelationType, constants.MergeType, constants.JsonbType:
			jdbcType = types.SrcJdbcTypePostgresql.Id()
		default:
			jdbcType = types.SrcJdbcTypeNone.Id()
		}
		unsDto.DataSrcID = jdbcType
	}
}
func newUnsFile(unsDto *types.CreateTopicDto) *dao.UnsNamespace {
	alias := unsDto.Alias
	instance := &dao.UnsNamespace{
		Id:          unsDto.Id,
		Alias:       alias,
		Name:        unsDto.Name,
		Fields:      unsDto.Fields,
		PathType:    unsDto.PathType,
		DataType:    unsDto.DataType,
		Description: unsDto.Description,
		DataPath:    unsDto.DataPath,
		DisplayName: unsDto.DisplayName,
		ParentAlias: unsDto.ParentAlias,
		Expression:  unsDto.Expression,
		Extend:      unsDto.Extend,
		ModelId:     unsDto.ModelId,
		WithFlags:   unsDto.WithFlags,
		MountType:   unsDto.MountType,
		MountSource: unsDto.MountSource,
	}

	if jdbcType := unsDto.DataSrcID; jdbcType != 0 {
		instance.DataSrcId = jdbcType
	}

	if unsDto.TableName != "" {
		instance.TableName_ = &unsDto.TableName
	}
	if unsDto.ModelAlias != nil {
		instance.ModelAlias = *unsDto.ModelAlias
	}
	if unsDto.PathType == constants.PathTypeFile && unsDto.ParentDataType != nil {
		instance.ParentDataType = unsDto.ParentDataType
	}

	if unsDto.Refers != nil {
		instance.Refers = unsDto.Refers
	}
	//instance.CalculationType = unsDto.CalculationType
	if len(unsDto.JsonFields) > 0 {
		protocol := unsDto.Protocol
		if protocol == nil {
			protocol = make(map[string]interface{})
			unsDto.Protocol = protocol
		}
		jsfStr, er := JsonUtil.ToJson(unsDto.JsonFields)
		if er == nil {
			protocol["jsf"] = jsfStr
		}
	}
	if protocol := unsDto.Protocol; protocol != nil && len(protocol) > 0 {
		if protocolType, exists := protocol["protocol"]; exists && protocolType != nil {
			instance.ProtocolType = base.V2p(fmt.Sprintf("%v", protocolType))
		}
		protocolBean := unsDto.ProtocolBean
		if protocolBean == nil {
			protocolBean = protocol
		}
		jsonProt, _ := JsonUtil.ToJson(protocolBean)
		instance.Protocol = base.V2p(jsonProt)
	}

	//if unsDto.ReadWriteMode != "" {
	//	instance.ReadWriteMode = unsDto.ReadWriteMode
	//}

	if unsDto.ExtendFieldUsed != nil {
		extTag := FieldFlags.GenerateFlag(unsDto.ExtendFieldUsed)
		instance.ExtendFieldFlags = &extTag
	}

	return instance
}
func getTemplate(topicDto *types.CreateTopicDto, existsUns func(string) *dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace) (template *dao.UnsNamespace, errMsg string) {
	modelId := topicDto.ModelId
	modelAlias := topicDto.ModelAlias
	var folderAlias *string

	if modelId != nil && *modelId != 0 {
		template = dbFiles[*modelId]
		if template != nil && template.PathType != 1 {
			errMsg = I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(template.PathType))),
				I18nUtils.GetMessage("uns.type.1"),
			)
		}
	} else if modelAlias != nil {
		template = existsUns(*modelAlias)
		if template == nil || base.P2v(template.Status) == LOGIC_REMOVED {
			errMsg = I18nUtils.GetMessage("uns.template.not.exists")
			template = nil
		} else if template.PathType != 1 {
			errMsg = I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(template.PathType))),
				I18nUtils.GetMessage("uns.type.1"),
			)
		}
	} else if folderAlias = topicDto.ParentAlias; folderAlias != nil {
		folder := existsUns(*folderAlias)
		if folder == nil {
			errMsg = I18nUtils.GetMessage("uns.folder.not.found") + ":alias=" + *folderAlias
		} else if pt := topicDto.PathType; folder.PathType != constants.PathTypeDir && (pt == constants.PathTypeFile || pt == constants.PathTypeDir) {
			errMsg = I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(folder.PathType))),
				I18nUtils.GetMessage("uns.type.0"),
			)
		} else if pt == constants.PathTypeTemplate && folder.PathType != constants.PathTypeTemplate {
			errMsg = I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(folder.PathType))),
				I18nUtils.GetMessage("uns.type.1"),
			)
		}
	}
	return template, errMsg
}
func setFieldsErr(unsDto *types.CreateTopicDto, errTipMap map[string]string, batchIndex string, instance *dao.UnsNamespace, template *dao.UnsNamespace) bool {
	insFs := unsDto.Fields
	jdbcType := types.SrcJdbcType(unsDto.DataSrcID)
	if len(insFs) == 0 && base.P2v(unsDto.DataType) == constants.JsonbType {
		insFs = []*types.FieldDefine{{Name: "json", Type: types.FieldTypeString}}
	}
	addSystemField := jdbcType != 0 && unsDto.PathType == constants.PathTypeFile && base.P2v(unsDto.DataType) != constants.AlarmRuleType

	if len(insFs) > 0 {
		tfd, err := FieldUtils.ProcessFieldDefines(jdbcType, insFs, true, addSystemField)
		if err != nil {
			errTipMap[batchIndex] = err.Error()
			return true
		}
		if tfd != nil {
			insFs = tfd.Fields
			instance.TableName_ = &tfd.TableName
		}
		instance.Fields = insFs
	} else {
		insFs = nil
	}

	if template != nil && template.Fields != nil {
		fields := template.Fields
		if len(insFs) > 0 {
			var checkError string
			if unsDto.PathType == constants.PathTypeFile && base.P2v(unsDto.DataType) != constants.JsonbType {
				checkError = checkInstanceFields(fields, insFs)
			}
			if checkError != "" {
				errTipMap[batchIndex] = checkError
				return true
			} else if instance.Fields == nil {
				tfd, err := FieldUtils.ProcessFieldDefines(jdbcType, fields, true, true)
				if err != nil {
					errTipMap[batchIndex] = err.Error()
					return true
				}
				if tfd != nil {
					insFs = tfd.Fields
					instance.TableName_ = &tfd.TableName
				}
				instance.Fields = insFs
			}
		} else if addSystemField {
			tfd, err := FieldUtils.ProcessFieldDefines(jdbcType, fields, true, true)
			if err != nil {
				errTipMap[batchIndex] = err.Error()
				return true
			}
			if tfd != nil {
				insFs = tfd.Fields
				instance.TableName_ = &tfd.TableName
			}
			unsDto.Fields = insFs
			instance.Fields = insFs
		}
	}

	if unsDto.PathType == constants.PathTypeFile && len(instance.Fields) == 0 {
		errTipMap[batchIndex] = I18nUtils.GetMessage("uns.field.empty")
		return true
	}

	return false
}
func (u *UnsAddService) trySetId(
	ctx context.Context,
	skipWhenExists bool,
	ct time.Time,
	unsDto *types.CreateTopicDto,
	existsUns func(string) *dao.UnsNamespace,
	dbFiles map[int64]*dao.UnsNamespace,
	addUpdate func(*dao.UnsNamespace),
	deleteFiles *[]*dao.UnsNamespace,
	errTipMap map[string]string) (po *dao.UnsNamespace, exists bool) {

	batchIndex := unsDto.GainBatchIndex()
	template, errMsg := getTemplate(unsDto, existsUns, dbFiles)
	if errMsg != "" {
		errTipMap[batchIndex] = errMsg
		return nil, false
	}

	dbPo := existsUns(unsDto.Alias)
	if dbPo != nil {
		if base.P2v(dbPo.Status) == OK && dbPo.PathType != unsDto.PathType {
			msg := I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(dbPo.PathType))),
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(unsDto.PathType))),
			)
			errTipMap[batchIndex] = msg
			return nil, base.P2v(dbPo.Status) == OK
		}
		unsDto.Id = dbPo.Id
	} else {
		unsDto.Id = common.NextId()
	}

	DB_EXISTS := dbPo != nil && base.P2v(dbPo.Status) == OK
	if !DB_EXISTS {
		setJdbcType(unsDto)
	}
	if DB_EXISTS && skipWhenExists {
		*unsDto = *UnsConverter.Po2Dto(dbPo)
		return dbPo, true
	}

	// 创建关系型文件, 不允许新增系统字段
	//if !DB_EXISTS && len(unsDto.Fields) > 0 &&
	//	unsDto.DataSrcID != 0 && types.SrcJdbcType(unsDto.DataSrcID).TypeCode() != constants.TimeSequenceType {
	//	for _, fd := range unsDto.Fields {
	//		if fd.IsSystemField() {
	//			errTipMap[batchIndex] = I18nUtils.GetMessage("uns.field.keyword", fd.Name) + ",alias=" + unsDto.Alias
	//			return nil
	//		}
	//	}
	//}

	newUns := newUnsFile(unsDto)
	if DB_EXISTS {
		tar := *dbPo
		_ = copier.CopyWithOption(&tar, newUns, copier.Option{IgnoreEmpty: true})
		if newUns.ParentAlias != nil && *newUns.ParentAlias == "" {
			tar.ParentAlias = nil
			tar.ParentId = nil
		}
		unsDto = UnsConverter.Po2Dto(&tar)
		newUns = &tar
		newUns.UpdateAt = ct

	}

	dataType := int16(0)
	if dt := unsDto.DataType; dt != nil {
		dataType = *dt
	}
	if (!DB_EXISTS && dataType != constants.CitingType) || len(unsDto.Fields) > 0 {
		if setFieldsErr(unsDto, errTipMap, batchIndex, newUns, template) {
			return nil, DB_EXISTS
		}
	}
	if len(unsDto.Fields) > 0 {
		newUns.NumberFields = base.V2p(unsDto.CountNumberFields())
	}
	if dataType == constants.CitingType && unsDto.Fields != nil {
		EMPTY := make([]*types.FieldDefine, 0)
		unsDto.Fields = EMPTY
		newUns.Fields = EMPTY
	}

	if DB_EXISTS {
		expression := newUns.Expression
		expChanged := expression != nil && (dbPo.Expression == nil || *expression != *dbPo.Expression)
		hasRefer := unsDto.Refers != nil || unsDto.ReferIds != nil

		checkFileFieldError := u.unsCalcService.CheckFileField(unsDto)
		if checkFileFieldError != "" {
			errTipMap[batchIndex] = checkFileFieldError
			return nil, DB_EXISTS
		}

		if expChanged || hasRefer {
			if hasRefer {
				err := u.unsCalcService.CheckRefers(unsDto)
				if err != "" {
					errTipMap[batchIndex] = err
					return nil, DB_EXISTS
				}
			}
			if expChanged {
				err := u.unsCalcService.CheckComplexExpression(unsDto)
				if err != "" {
					errTipMap[batchIndex] = err
					return nil, DB_EXISTS
				}
			}
		}
		newUns.Refers = unsDto.Refers
		newUns.Expression = unsDto.Expression
		var fieldsChanged = !base.EqualsF(newUns.Fields, dbPo.Fields, func(a, b *types.FieldDefine) bool {
			return a.Name == b.Name && a.Type == b.Type
		})
		unsDto.FieldsChanged = fieldsChanged
		isFile, isTemplate := newUns.PathType == constants.PathTypeFile, newUns.PathType == constants.PathTypeTemplate
		if fieldsChanged && (isFile || isTemplate) {
			oldFields := normalFields(dbPo.Fields)
			inputsMap := base.MapArrayToMap[*types.FieldDefine, string, *types.FieldDefine](newUns.Fields, func(e *types.FieldDefine) (ok bool, k string, v *types.FieldDefine) {
				return !e.IsSystemField(), e.Name, e
			})
			// 筛选出删除的属性集合
			delFields := base.Filter(oldFields, func(e *types.FieldDefine) bool {
				return !base.MapContainsKey(inputsMap, e.Name)
			})

			var affected []*dao.UnsNamespace
			if isFile {
				affected = []*dao.UnsNamespace{newUns}
			} else if isTemplate {
				var err error
				db := dao.GetDb(ctx)
				affected, err = u.unsMapper.ListByTemplateId(db, newUns.Id, nil)
				if err != nil {
					errTipMap[batchIndex] = err.Error()
					return nil, DB_EXISTS
				}
				if len(affected) > 0 {
					for _, f := range affected {
						tdbFs, er := FieldUtils.ProcessFieldDefines(types.SrcJdbcType(f.DataSrcId), newUns.Fields, true, true)
						if er == nil && tdbFs != nil {
							if base.P2v(f.DataType) == constants.JsonbType {
								jsfStr, _ := JsonUtil.ToJson(newUns.Fields)
								pm := f.GetProtocolMap()
								if pm == nil {
									pm = make(map[string]interface{})
								}
								pm["jsf"] = jsfStr
								jsfStr, _ = JsonUtil.ToJson(pm)
								f.Protocol = &jsfStr
							} else {
								f.Fields = tdbFs.Fields
							}
							addUpdate(f)
						}
						dbFiles[f.Id] = f
					}
				}
			}
			if len(delFields) > 0 && len(affected) > 0 {
				affectedDeleteFiles := u.unsCalcService.detectReferencedCalcInstance(ctx, affected, newUns.Path, delFields)
				if len(affectedDeleteFiles) > 0 {
					*deleteFiles = append(*deleteFiles, affectedDeleteFiles...)
				}
			}
		}
	} else {
		newUns.Name = strings.TrimSpace(newUns.Name)
		if len(newUns.Name) == 0 {
			errTipMap[batchIndex] = I18nUtils.GetMessage("uns.name.empty")
			return nil, DB_EXISTS
		}
		if unsDto.WithFlags == nil {
			flag := generateFlag(unsDto.AddFlow, unsDto.Save2Db, unsDto.AddDashBoard,
				unsDto.RetainTableWhenDeleteInstance, unsDto.SubscribeEnable, unsDto.AccessLevel)
			newUns.WithFlags = &flag
		}
		//if newUns.ReadWriteMode == "" {
		//	newUns.ReadWriteMode = FileReadWriteModeReadOnly.Mode
		//	unsDto.ReadWriteMode = newUns.ReadWriteMode
		//}
		if newUns.ExtendFieldFlags == nil {
			extFlag := FieldFlags.GenerateFlag(unsDto.ExtendFieldUsed)
			newUns.ExtendFieldFlags = &extFlag
		}

		var err string
		err = u.unsCalcService.CheckFileField(unsDto)
		if err == "" {
			err = u.unsCalcService.CheckRefers(unsDto)
		}
		if err == "" {
			err = u.unsCalcService.CheckComplexExpression(unsDto)
		}
		if err != "" {
			errTipMap[batchIndex] = err
			newUns = nil
		} else {
			newUns.CreateAt = ct
			newUns.Refers = unsDto.Refers
			newUns.Expression = unsDto.Expression
		}
	}

	if newUns != nil {
		unsDto.Status = 1
		newUns.Status = &OK
	}
	return newUns, DB_EXISTS
}
func normalFields(fs []*types.FieldDefine) []*types.FieldDefine {
	return base.Filter(fs, func(e *types.FieldDefine) bool {
		return !e.IsSystemField()
	})
}
func generateFlag(addFlow, saveToDB, addDashBoard, retainTableWhenDeleteInstance, subscribeEnable *bool, accessLevel string) int32 {
	flags := int32(0)

	if is(addFlow) {
		flags |= constants.UnsFlagWithFlow
	}
	if is(saveToDB) {
		flags |= constants.UnsFlagWithSave2DB
	}
	if is(addDashBoard) {
		flags |= constants.UnsFlagWithDashboard
	}
	if is(retainTableWhenDeleteInstance) {
		flags |= constants.UnsFlagRetainTableWhenDelInstance
	}
	if accessLevel == constants.AccessLevelReadOnly {
		flags |= constants.UnsFlagAccessLevelReadOnly
	}
	if accessLevel == constants.AccessLevelReadWrite {
		flags |= constants.UnsFlagAccessLevelReadWrite
	}
	if is(subscribeEnable) {
		flags |= constants.UnsFlagWithSubscribeEnable
	}

	return flags
}
func is(b *bool) bool {
	return b != nil && *b
}

type unsDtoTreeNodes struct {
	uns []*dao.UnsNamespace
}

func (u *unsDtoTreeNodes) Size() int {
	return len(u.uns)
}
func (u *unsDtoTreeNodes) Visit(visitor func(uns *dao.UnsNamespace)) {
	for _, node := range u.uns {
		visitor(node)
	}
}

type Siblings struct {
	names map[string][]*dao.UnsNamespace
}

func newSiblings() *Siblings {
	return &Siblings{names: make(map[string][]*dao.UnsNamespace, 32)}
}
func (s *Siblings) add(uns *dao.UnsNamespace) {
	s.names[uns.Name] = append(s.names[uns.Name], uns)
}

const _NullParentId = int64(-1)

func scanChangedNodes(addFiles map[int64]*dao.UnsNamespace, existsUns map[string]*dao.UnsNamespace, siblings map[int64]*Siblings, changedSubTree *[]*dao.UnsNamespace) {
	for _, bo := range addFiles {
		alias := bo.Alias
		dbo := existsUns[alias]
		parentAlias := bo.ParentAlias
		var scanSiblings = false
		if dbo == nil {
			scanSiblings = true
		} else if !eqStrP(parentAlias, dbo.ParentAlias) || bo.Name != dbo.Name {
			scanSiblings = true
			if bo.PathType == constants.PathTypeDir {
				*changedSubTree = append(*changedSubTree, dbo)
			}
		}
		if scanSiblings {
			parentId := base.P2vWithDefault(bo.ParentId, _NullParentId)
			sib, has := siblings[parentId]
			if !has {
				sib = newSiblings()
				siblings[parentId] = sib
			}
			sib.add(bo)
		}
	}
}
func (u *UnsAddService) tryAddLayRecOrPathChangedChildren(ctx context.Context,
	addFiles map[int64]*dao.UnsNamespace,
	dbFiles map[int64]*dao.UnsNamespace,
	existsUns map[string]*dao.UnsNamespace,
) error {
	changedSubTree := make([]*dao.UnsNamespace, 0, 64)
	parentIdSet := make(map[int64]*Siblings, 32)
	scanChangedNodes(addFiles, existsUns, parentIdSet, &changedSubTree)
	sizeTree, sizeSiblings := len(changedSubTree), len(parentIdSet)
	if sizeTree+sizeSiblings == 0 {
		return nil
	}
	db := dao.GetDb(ctx)
	if sizeTree > 0 {
		topNodes := changedSubTree
		if len(topNodes) > 1 {
			topNodes = LeastTopNodeUtil.GetLeastTopNodes(&unsDtoTreeNodes{uns: topNodes})
		}
		for _, po := range topNodes {
			children, er := u.unsMapper.ListSubTree(db, po.LayRec)
			if er != nil {
				return er
			}
			if len(children) > 0 {
				for _, unsPo := range children {
					unsPo.LayRec = "" //重置，等着重新计算
					dbFiles[unsPo.Id] = unsPo
					existsUns[unsPo.Alias] = unsPo
				}
			}
		}
		for _, unsPo := range changedSubTree {
			unsPo.LayRec = "" //重置，等着重新计算
		}
	}
	if sizeSiblings > 0 {
		var siblings = base.FilterAndFlatMap(base.MapValues(parentIdSet), func(sib *Siblings) (vs []*dao.UnsNamespace, ok bool) {
			vs = make([]*dao.UnsNamespace, len(sib.names))
			i := 0
			for name, cs := range sib.names {
				vs[i] = &dao.UnsNamespace{ParentId: cs[0].ParentId, Name: name, Alias: cs[0].Alias}
				i++
			}
			return vs, true
		})
		for _, partSiblings := range base.Partition(siblings, 1000) {
			countMap, er := u.unsMapper.CountByParentAliasAndNames(db, partSiblings)
			if er != nil {
				return er
			}
			if len(countMap) > 0 {
				for _, cm := range countMap {
					parentId := base.P2vWithDefault(cm.ParentId, _NullParentId)
					sib := parentIdSet[parentId]
					if sib != nil && len(sib.names) > 0 {
						sameNameSiblings := sib.names[cm.Name]
						if len(sameNameSiblings) > 0 {
							for _, uns := range sameNameSiblings {
								uns.LayRec = ""
								uns.CountExistsSiblings = cm.CountExistsSiblings
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func aliasToId(addFiles map[int64]*dao.UnsNamespace, aliasMap func(string) *dao.UnsNamespace, pathMap map[int16]map[string]*types.CreateTopicDto) {
	for _, file := range addFiles {
		if modelAlias := file.ModelAlias; modelAlias != "" {
			if model := aliasMap(modelAlias); model != nil {
				file.ModelId = &model.Id
			}
		}
		if parentAlias := file.ParentAlias; parentAlias != nil {
			if parent := aliasMap(*parentAlias); parent != nil {
				file.ParentId = &parent.Id
			}
		}
		if dtoMap, has := pathMap[file.PathType]; has {
			DTO := dtoMap[file.Alias]
			if DTO != nil {
				DTO.Id = file.Id
				DTO.ModelId = file.ModelId
				DTO.ParentId = file.ParentId
			}
		}
	}
}

func (u *UnsAddService) listUnsByAliasAndIds(ctx context.Context, alias []string, ids []int64, dbFiles map[int64]*dao.UnsNamespace, aliasMap map[string]*dao.UnsNamespace) (er error) {
	db := dao.GetDb(ctx)
	for _, aliasList := range base.Partition(alias, constants.SQLBatchSize) {
		unsPos, er := u.unsMapper.AllByAlias(db, aliasList)
		if er != nil {
			return er
		}
		addDbPo(unsPos, dbFiles, aliasMap)
	}
	if len(ids) > 0 {
		for _, idList := range base.Partition(ids, constants.SQLBatchSize) {
			unsPos, er := u.unsMapper.AllByIds(db, idList)
			if er != nil {
				return er
			}
			addDbPo(unsPos, dbFiles, aliasMap)
		}
	}
	return
}

func getEventStatusCallback(statusConsumer func(status *common.RunningStatus)) event.EventStatusAware {
	if statusConsumer == nil {
		return nil
	}
	return newWrappedEventStatusAware(statusConsumer)
}

type wrappedEventStatusAware struct {
	t0             int64
	statusConsumer func(status *common.RunningStatus)
}

var _startMsg, _endMsg, _errMsg string
var once sync.Once

func newWrappedEventStatusAware(statusConsumer func(status *common.RunningStatus)) event.EventStatusAware {
	once.Do(func() {
		_startMsg = I18nUtils.GetMessage("uns.create.status.running")
		_endMsg = I18nUtils.GetMessage("uns.create.status.finished")
		_errMsg = I18nUtils.GetMessage("uns.create.status.error")
	})
	return &wrappedEventStatusAware{statusConsumer: statusConsumer}
}
func (w *wrappedEventStatusAware) BeforeEvent(N int, i int, listenerName string) {
	progress := 0.0
	if i > 1 && N > 0 {
		progress = float64(int((1000.0 * (float64(i) - 1) / float64(N)))) / 10.0
	}
	w.statusConsumer(common.NewRunningStatusWithProgress(N, i, listenerName, _startMsg).SetProgress(progress))
	w.t0 = time.Now().UnixMilli()
}
func (w *wrappedEventStatusAware) AfterEvent(N int, i int, listenerName string, err error) {
	msg := _endMsg
	code := 200
	if err != nil {
		code = 500
		msg = _errMsg + err.Error()
	}
	w.statusConsumer(common.NewRunningStatusWithProgress(N, i, listenerName, msg).
		SetSpendMills(time.Now().UnixMilli() - w.t0).SetCode(code))
}
