package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/LeastTopNodeUtil"
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	"backend/internal/common/enums"
	"backend/internal/common/event"
	"backend/internal/common/utils/FieldUtils"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/share/base"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

func checkInstanceFields(modelFields []*dto.FieldDefine, insFields []*dto.FieldDefine) string {
	if modelFields == nil {
		modelFields = make([]*dto.FieldDefine, 0)
	}
	if insFields == nil {
		insFields = make([]*dto.FieldDefine, 0)
	}

	insMap := make(map[string]*dto.FieldDefine, len(insFields))
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

func initParamsUns(topicDtos []*dto.CreateTopicDto, errTipMap map[string]string, paramFiles map[string]*dto.CreateTopicDto) map[string]*dto.CreateTopicDto {
	paramFolders := make(map[string]*dto.CreateTopicDto, 2+len(topicDtos)/2)
	for _, topicDto := range topicDtos {
		checkTopicDto(errTipMap, paramFolders, paramFiles, topicDto)
	}
	return paramFolders
}
func addAlias(bos []*dto.CreateTopicDto, aliasSet map[string]bool, ids map[int64]bool) {
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
		if modelAlias := unsDto.ModelAlias; modelAlias != "" {
			aliasSet[modelAlias] = true
		}

		// 添加parentAlias
		if folderAlias := unsDto.ParentAlias; folderAlias != nil {
			aliasSet[*folderAlias] = true
		}

		// 处理referIds
		if referIds := unsDto.ReferIDs; len(referIds) > 0 {
			// 添加所有referIds到ids集合
			for _, id := range referIds {
				ids[id] = true
			}

			// 如果refers为空，根据referIds创建InstanceField数组
			if len(unsDto.Refers) == 0 {
				refers := make([]*dto.InstanceField, len(referIds))
				for i, id := range referIds {
					refers[i] = &dto.InstanceField{ID: id, Alias: ""}
				}
				unsDto.Refers = refers
			}
		}

		// 添加各种ID到ids集合
		if unsId := unsDto.ID; unsId != 0 {
			ids[unsId] = true
		}
		if pid := unsDto.ParentID; pid != nil {
			ids[*pid] = true
		}
		if mid := unsDto.ModelID; mid != nil {
			ids[*mid] = true
		}

		// 处理refers中的字段
		if refers := unsDto.Refers; len(refers) > 0 {
			for _, field := range refers {
				if id := field.ID; id != 0 {
					ids[id] = true
				}
				if alias := field.Alias; alias != "" {
					aliasSet[alias] = true
				}
			}
		}
	}
}

func scanChangedNodes(files []*dto.CreateTopicDto, existsUns map[string]*dao.UnsNamespace, siblings map[string]*Siblings, changedSubTree *[]*dao.UnsNamespace) {
	for _, bo := range files {
		alias := bo.Alias
		dbo := existsUns[alias]
		parentAlias := bo.ParentAlias
		if dbo == nil {
		} else if bo.PathType == constants.PathTypeDir && (!eqStrP(parentAlias, dbo.ParentAlias) || bo.Name != dbo.Name) {
			*changedSubTree = append(*changedSubTree, dbo)
		}
		pAlias := ""
		if parentAlias != nil {
			pAlias = *parentAlias
		}
		sib, has := siblings[pAlias]
		if !has {
			sib = newSiblings()
			siblings[pAlias] = sib
		}
		sib.add(bo)
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
func tryFillIdOrAlias(paramFiles map[string]*dto.CreateTopicDto, existsUns map[string]*dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace, errTipMap map[string]string) {
	for key, topicDto := range paramFiles {
		// 处理主对象的ID和Alias
		id := topicDto.ID
		alias := topicDto.Alias
		if id != 0 && alias == "" {
			if po, exists := dbFiles[id]; exists {
				topicDto.Alias = po.Alias
			}
		} else if id == 0 && alias != "" {
			if po, exists := existsUns[alias]; exists {
				topicDto.ID = po.ID
			}
		}

		// 处理父级ID和ParentAlias
		pid := topicDto.ParentID
		parentAlias := topicDto.ParentAlias
		if pid == nil && parentAlias != nil {
			if parent, exists := existsUns[*parentAlias]; exists {
				topicDto.ParentID = &parent.ID
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
				refID := field.ID
				refAlias := field.Alias
				if refID == 0 && refAlias != "" {
					if refPo, exists := existsUns[refAlias]; exists {
						field.ID = refPo.ID
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

// 常量定义
const (
	LOGIC_REMOVED = 0
	OK            = 1
)

// 临时处理PO
func putTemp(dbFiles map[int64]*dao.UnsNamespace, aliasMap map[string]*dao.UnsNamespace, po *dao.UnsNamespace) {
	if po.Status == LOGIC_REMOVED {
		// 创建新的PO对象，只保留ID和Alias
		newPo := &dao.UnsNamespace{
			ID:     po.ID,
			Alias:  po.Alias,
			Status: LOGIC_REMOVED,
		}
		// 使用新对象替换原对象
		po = newPo
	}

	// 添加到映射
	aliasMap[po.Alias] = po
	dbFiles[po.ID] = po
}

func checkTopicDto(errTipMap map[string]string,
	paramFolders map[string]*dto.CreateTopicDto,
	paramFiles map[string]*dto.CreateTopicDto,
	d *dto.CreateTopicDto) {
	pathType := d.PathType
	if pathType == constants.PathTypeDir {
		d.DataType = nil
	}

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
	if base.MapContainsKey(paramFolders, alias) || base.MapContainsKey(paramFiles, alias) {
		msg := I18nUtils.GetMessage("uns.alias.duplicate")
		errTipMap[batchIndex] = msg
		return
	}

	if pathType == constants.PathTypeDir { // 当前是文件夹
		d.DataType = nil
		paramFolders[alias] = d
	} else if pathType == constants.PathTypeFile { // 当前是文件
		dataType := d.DataType
		if dataType == nil {
			msg := I18nUtils.GetMessage("uns.file.dataType.empty")
			errTipMap[batchIndex] = msg
			return
		} else if !constants.IsValidDataType(*dataType) {
			msg := I18nUtils.GetMessage(fmt.Sprintf("uns.file.dataType.invalid", *dataType))
			errTipMap[batchIndex] = msg
			return
		}

		fields := d.Fields
		if len(fields) == 0 && *dataType == constants.MergeType {
			mergeField := &dto.FieldDefine{
				Name:   "data_json",
				Type:   enums.FieldTypeString,
				MaxLen: 512 * 1024, // 聚合的字段总长度限制改大，不能超过mqtt消息长度限制
			}
			fields = []*dto.FieldDefine{mergeField}
			d.Fields = fields
		}
		paramFiles[alias] = d
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

func setJdbcType(unsDto *dto.CreateTopicDto) {
	dataType := unsDto.DataType
	jdbcType := unsDto.DataSrcID
	if jdbcType == 0 && dataType != nil && unsDto.PathType == constants.PathTypeFile {
		switch *dataType {
		case constants.CalculationHistType, constants.CalculationRealType, constants.TimeSequenceType:
			jdbcType = common.SrcJdbcTypeTimeScaleDB
		case constants.AlarmRuleType, constants.RelationType, constants.MergeType:
			jdbcType = common.SrcJdbcTypePostgresql
		default:
			jdbcType = common.SrcJdbcTypeNone
		}
		unsDto.DataSrcID = jdbcType
	}
}
func newUnsFile(unsDto *dto.CreateTopicDto) *dao.UnsNamespace {
	alias := unsDto.Alias
	instance := &dao.UnsNamespace{
		ID:                  unsDto.ID,
		Alias:               alias,
		Name:                unsDto.Name,
		PathType:            unsDto.PathType,
		Description:         unsDto.Description,
		CountExistsSiblings: unsDto.CountExistsSiblings,
	}

	if jdbcType := unsDto.DataSrcID; jdbcType != 0 {
		instance.DataSrcID = jdbcType.Id()
	}

	instance.DisplayName = unsDto.DisplayName
	instance.DataPath = unsDto.DataPath
	instance.Alias = alias
	instance.Name = unsDto.Name
	instance.ParentAlias = unsDto.ParentAlias
	instance.TableName_ = unsDto.TableName
	instance.WithFlags = unsDto.Flags
	instance.ModelID = unsDto.ModelID
	instance.ModelAlias = unsDto.ModelAlias

	if unsDto.PathType == constants.PathTypeFile {
		if dataType := unsDto.DataType; dataType != nil {
			instance.DataType = dataType
		}
		if len(unsDto.Fields) > 0 {
			instance.NumberFields = unsDto.CountNumberFields()
		}
	}

	instance.Extend = unsDto.Extend

	if unsDto.Refers != nil {
		instance.Refers = unsDto.Refers
	}

	instance.Expression = unsDto.Expression
	//instance.CalculationType = unsDto.CalculationType

	if protocol := unsDto.Protocol; protocol != nil && len(protocol) > 0 {
		if protocolType, exists := protocol["protocol"]; exists && protocolType != nil {
			instance.ProtocolType = fmt.Sprintf("%v", protocolType)
		}
		protocolBean := unsDto.ProtocolBean
		if protocolBean == nil {
			protocolBean = protocol
		}
		instance.Protocol, _ = JsonUtil.ToJson(protocolBean)
	}

	instance.MountType = unsDto.MountType
	instance.MountSource = unsDto.MountSource

	//if unsDto.ReadWriteMode != "" {
	//	instance.ReadWriteMode = unsDto.ReadWriteMode
	//}

	if unsDto.ExtendFieldUsed != nil {
		extTag := FieldUtils.GenerateFlag(unsDto.ExtendFieldUsed)
		instance.ExtendFieldFlags = &extTag
	}

	return instance
}
func getTemplate(topicDto *dto.CreateTopicDto, existsUns func(string) *dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace) (template *dao.UnsNamespace, errMsg string) {
	modelId := topicDto.ModelID
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
	} else if modelAlias != "" {
		template = existsUns(modelAlias)
		if template != nil && template.PathType != 1 {
			errMsg = I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(template.PathType))),
				I18nUtils.GetMessage("uns.type.1"),
			)
		} else if template == nil {
			errMsg = I18nUtils.GetMessage("uns.template.not.exists")
		}
	} else if folderAlias = topicDto.ParentAlias; folderAlias != nil {
		folder := existsUns(*folderAlias)
		if folder == nil {
			errMsg = I18nUtils.GetMessage("uns.folder.not.found")
		} else if folder.PathType != constants.PathTypeDir {
			errMsg = I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(folder.PathType))),
				I18nUtils.GetMessage("uns.type.0"),
			)
		}
	}
	return template, errMsg
}
func setFieldsErr(unsDto *dto.CreateTopicDto, errTipMap map[string]string, batchIndex string, instance *dao.UnsNamespace, template *dao.UnsNamespace) bool {
	insFs := unsDto.Fields
	jdbcType := unsDto.DataSrcID

	addSystemField := jdbcType != 0 && unsDto.PathType == constants.PathTypeFile && (unsDto.DataType != nil && *unsDto.DataType != constants.AlarmRuleType)

	if len(insFs) > 0 {
		tfd, err := FieldUtils.ProcessFieldDefines(jdbcType, insFs, true, addSystemField)
		if err != nil {
			errTipMap[batchIndex] = err.Error()
			return true
		}
		insFs = tfd.Fields
		instance.TableName_ = tfd.TableName
		instance.Fields = insFs
	} else {
		insFs = nil
	}

	if template != nil && template.Fields != nil {
		fields := template.Fields
		if len(insFs) > 0 {
			var checkError string
			if unsDto.PathType == constants.PathTypeFile {
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
				insFs = tfd.Fields
				instance.TableName_ = tfd.TableName
				instance.Fields = insFs
			}
		} else if addSystemField {
			tfd, err := FieldUtils.ProcessFieldDefines(jdbcType, fields, true, true)
			if err != nil {
				errTipMap[batchIndex] = err.Error()
				return true
			}
			insFs = tfd.Fields
			instance.TableName_ = tfd.TableName
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
func (u *UnsAddService) trySetId(ct time.Time, unsDto *dto.CreateTopicDto, existsUns func(string) *dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace, errTipMap map[string]string) *dao.UnsNamespace {
	batchIndex := unsDto.GainBatchIndex()
	template, errMsg := getTemplate(unsDto, existsUns, dbFiles)
	if errMsg != "" {
		errTipMap[batchIndex] = errMsg
		return nil
	}
	setJdbcType(unsDto)

	dbPo := existsUns(unsDto.Alias)
	if dbPo != nil {
		if dbPo.PathType != unsDto.PathType {
			msg := I18nUtils.GetMessage("uns.alias.has.exist.type",
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(dbPo.PathType))),
				I18nUtils.GetMessage("uns.type."+strconv.Itoa(int(unsDto.PathType))),
			)
			errTipMap[batchIndex] = msg
			return nil
		}
		unsDto.ID = dbPo.ID
	} else {
		unsDto.ID = common.NextId()
	}

	DB_EXISTS := dbPo != nil && dbPo.Status == OK

	// 创建关系型文件, 不允许新增系统字段
	if !DB_EXISTS && len(unsDto.Fields) > 0 &&
		unsDto.DataSrcID != 0 && unsDto.DataSrcID.TypeCode() != constants.TimeSequenceType {
		for _, fd := range unsDto.Fields {
			if fd.IsSystemField() {
				errTipMap[batchIndex] = I18nUtils.GetMessage("uns.field.keyword", fd.Name)
				return nil
			}
		}
	}

	newUns := newUnsFile(unsDto)
	dataType := int16(0)
	if dt := unsDto.DataType; dt != nil {
		dataType = *dt
	}
	if (!DB_EXISTS && dataType != constants.CitingType) || len(unsDto.Fields) > 0 {
		if setFieldsErr(unsDto, errTipMap, batchIndex, newUns, template) {
			return nil
		}
	}

	if dataType == constants.CitingType && unsDto.Fields != nil {
		EMPTY := make([]*dto.FieldDefine, 0)
		unsDto.Fields = EMPTY
		newUns.Fields = EMPTY
	}

	if DB_EXISTS {
		tar := *dbPo
		copier.CopyWithOption(&tar, newUns, copier.Option{IgnoreEmpty: true})
		expression := newUns.Expression
		expChanged := expression != "" && expression != dbPo.Expression
		hasRefer := unsDto.Refers != nil || unsDto.ReferIDs != nil

		checkFileFieldError := u.unsCalcService.CheckFileField(unsDto)
		if checkFileFieldError != "" {
			errTipMap[batchIndex] = checkFileFieldError
			return nil
		}

		if expChanged || hasRefer {
			unsDto = UnsConverter.Po2Dto(&tar)
			if hasRefer {
				err := u.unsCalcService.CheckRefers(unsDto)
				if err != "" {
					errTipMap[batchIndex] = err
					return nil
				}
			}
			if expChanged {
				err := u.unsCalcService.CheckComplexExpression(unsDto)
				if err != "" {
					errTipMap[batchIndex] = err
					return nil
				}
			}
		}
		tar.UpdateAt = ct
		newUns = &tar
		newUns.Refers = unsDto.Refers
		newUns.Expression = unsDto.Expression
	} else {
		if unsDto.Flags == nil {
			flag := generateFlag(unsDto.AddFlow, unsDto.Save2DB, unsDto.AddDashBoard,
				unsDto.RetainTableWhenDeleteInstance, unsDto.SubscribeEnable, unsDto.AccessLevel)
			newUns.WithFlags = &flag
		}
		//if newUns.ReadWriteMode == "" {
		//	newUns.ReadWriteMode = FileReadWriteModeReadOnly.Mode
		//	unsDto.ReadWriteMode = newUns.ReadWriteMode
		//}
		if newUns.ExtendFieldFlags == nil {
			extFlag := FieldUtils.GenerateFlag(unsDto.ExtendFieldUsed)
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
		newUns.Status = 1
	}
	return newUns
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
	names map[string][]*dto.CreateTopicDto
}

func newSiblings() *Siblings {
	return &Siblings{names: make(map[string][]*dto.CreateTopicDto, 32)}
}
func (s *Siblings) add(uns *dto.CreateTopicDto) {
	s.names[uns.Name] = append(s.names[uns.Name], uns)
}
func (u *UnsAddService) tryAddLayRecOrPathChangedChildren(ctx context.Context, paramFolders []*dto.CreateTopicDto, paramFiles []*dto.CreateTopicDto, existsUns map[string]*dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace) error {
	changedSubTree := make([]*dao.UnsNamespace, 0, 64)
	parentAliasSet := make(map[string]*Siblings, 32)
	scanChangedNodes(paramFolders, existsUns, parentAliasSet, &changedSubTree)
	scanChangedNodes(paramFiles, existsUns, parentAliasSet, &changedSubTree)
	sizeTree, sizeSiblings := len(changedSubTree), len(parentAliasSet)
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
					dbFiles[unsPo.ID] = unsPo
					existsUns[unsPo.Alias] = unsPo
				}
			}
		}
		for _, unsPo := range changedSubTree {
			unsPo.LayRec = "" //重置，等着重新计算
		}
	}
	if sizeSiblings > 0 {
		var siblings = base.FilterAndFlatMap(base.MapValues(parentAliasSet), func(sib *Siblings) (vs []*dao.UnsNamespace, ok bool) {
			vs = make([]*dao.UnsNamespace, len(sib.names))
			i := 0
			for name, cs := range sib.names {
				vs[i] = &dao.UnsNamespace{ParentAlias: cs[0].ParentAlias, Name: name}
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
					parentAlias := ""
					if pa := cm.ParentAlias; pa != nil {
						parentAlias = *pa
					}
					sib := parentAliasSet[parentAlias]
					if sib != nil && len(sib.names) > 0 {
						sameNameSiblings := sib.names[cm.Name]
						if len(sameNameSiblings) > 0 {
							for _, uns := range sameNameSiblings {
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

func aliasToId(addFiles map[int64]*dao.UnsNamespace, aliasMap func(string) *dao.UnsNamespace) {
	for _, file := range addFiles {
		if modelAlias := file.ModelAlias; modelAlias != "" {
			if model := aliasMap(modelAlias); model != nil {
				file.ModelID = &model.ID
			}
		}
		if parentAlias := file.ParentAlias; parentAlias != nil {
			if parent := aliasMap(*parentAlias); parent != nil {
				file.ParentID = &parent.ID
			}
		}
	}
}

type unsLevel struct {
	uns      []*dto.CreateTopicDto
	levelMap map[string]int
}

func (x *unsLevel) Len() int { return len(x.uns) }
func (x *unsLevel) Less(i, j int) bool {
	a, b := x.levelMap[x.uns[i].Alias], x.levelMap[x.uns[j].Alias]
	return a < b
}
func (x *unsLevel) Swap(i, j int) { x.uns[i], x.uns[j] = x.uns[j], x.uns[i] }

func (u *UnsAddService) listUnsByAliasAndIds(ctx context.Context, alias []string, ids []int64, dbFiles map[int64]*dao.UnsNamespace) (aliasMap map[string]*dao.UnsNamespace, er error) {
	aliasMap = make(map[string]*dao.UnsNamespace, len(alias)+len(ids))
	db := dao.GetDb(ctx)
	for _, aliasList := range base.Partition(alias, constants.SQLBatchSize) {
		unsPos, er := u.unsMapper.ListByAlias(db, aliasList)
		if er != nil {
			return nil, er
		}
		addDbPo(unsPos, dbFiles, aliasMap)
	}
	if len(ids) > 0 {
		for _, idList := range base.Partition(ids, constants.SQLBatchSize) {
			unsPos, er := u.unsMapper.ListByIds(db, idList)
			if er != nil {
				return nil, er
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
	w.statusConsumer(common.NewRunningStatusWithProgress(N+1, i, listenerName, _startMsg).SetProgress(progress))
	w.t0 = time.Now().UnixMilli()
}
func (w *wrappedEventStatusAware) AfterEvent(N int, i int, listenerName string, err error) {
	msg := _endMsg
	code := 0
	if err != nil {
		code = 500
		msg = _errMsg + err.Error()
	}
	w.statusConsumer(common.NewRunningStatusWithProgress(N+1, i, listenerName, msg).
		SetSpendMills(time.Now().UnixMilli() - w.t0).SetCode(code))
}
