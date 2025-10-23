package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	"backend/internal/common/event"
	"backend/internal/logic/supos/uns/label/service"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsAddService struct {
	log             logx.Logger
	unsMapper       dao.UnsNamespaceRepo
	labelRefMapper  dao.UnsLabelRefRepo
	unsLabelService *service.UnsLabelService
	unsCalcService  UnsCalcService
}

func init() {
	spring.RegisterLazy[*UnsAddService](func() *UnsAddService {
		return &UnsAddService{
			log:             logx.WithContext(context.Background()),
			unsLabelService: spring.GetBean[*service.UnsLabelService](),
		}
	})
}
func (u *UnsAddService) CreateModelAndInstancesInner(ctx context.Context, args bo.CreateModelInstancesArgs) (errTipMap map[string]string) {
	topicDtos := args.Topics
	errTipMap = make(map[string]string, len(topicDtos))
	paramFiles := make(map[string]*dto.CreateTopicDto)
	paramFolders := initParamsUns(topicDtos, errTipMap, paramFiles)
	if len(paramFolders) == 0 && len(paramFiles) == 0 {
		u.log.Info("不存在任何文件夹或文件, 无法继续保存")
		return errTipMap
	}
	dbFiles := make(map[int64]*dao.UnsNamespace)
	var existsUns map[string]*dao.UnsNamespace
	{
		db := dao.GetDb(ctx)
		ctx = dao.SetDb(ctx, db)
	}
	{
		ids := make(map[int64]bool)
		aliasSet := make(map[string]bool)
		addAlias(base.MapValues(paramFolders), aliasSet, ids)
		addAlias(base.MapValues(paramFiles), aliasSet, ids)
		var err error
		existsUns, err = u.listUnsByAliasAndIds(ctx, base.MapKeys(aliasSet), base.MapKeys(ids), dbFiles)
		if err != nil {
			u.log.Error("QueryUnsError:", err)
			errTipMap["0"] = "QueryUnsError:" + err.Error()
			return errTipMap
		}
	}
	tryFillIdOrAlias(paramFiles, existsUns, dbFiles, errTipMap)
	addFiles := make(map[int64]*dao.UnsNamespace)
	aliasMap := make(map[string]*dao.UnsNamespace)
	folders := base.MapValues(paramFolders)
	if len(folders) > 1 {
		reverseGraph := base.BuildReverseGraph(folders, func(t *dto.CreateTopicDto) string {
			return t.Alias
		}, func(t *dto.CreateTopicDto) string {
			return t.ParentAlias
		})
		levelMap := base.CalculateLevels(reverseGraph)
		sort.Sort(&unsLevel{uns: folders, levelMap: levelMap})
	}
	// 找出 parentAlias 或 name 有修改的最高层目录，后面需要获取它的整个子树，为更新 layRec 做准备
	err := u.tryAddLayRecOrPathChangedChildren(ctx, folders, base.MapValues(paramFiles), existsUns, dbFiles)
	if err != nil {
		errTipMap["0"] = err.Error()
		return errTipMap
	}
	createTime := time.Now()
	allUns := func(alias string) *dao.UnsNamespace {
		unsPo := aliasMap[alias]
		if unsPo == nil {
			unsPo = existsUns[alias]
		}
		return unsPo
	}
	for _, folder := range folders {
		po := u.trySetId(createTime, folder, allUns, dbFiles, errTipMap)
		if po != nil {
			addFiles[po.ID] = po
			aliasMap[po.Alias] = po
		}
	}

	unsPoLabels := make(map[int64]*bo.UnsPoLabels, len(paramFiles))
	for _, unsFile := range paramFiles {
		po := u.trySetId(createTime, unsFile, allUns, dbFiles, errTipMap)
		if po != nil {
			addFiles[po.ID] = po
			aliasMap[po.Alias] = po
			if unsFile.LabelNames != nil {
				_, exists := dbFiles[po.ID]
				unsPoLabels[po.ID] = bo.NewUnsPoLabels(po, exists, unsFile.LabelNames)
			}
		}
	}
	//TODO 计算，引用，聚合等类型的 校验和处理
	aliasToId(addFiles, allUns)
	rs := setLayRecAndPath(createTime, addFiles, dbFiles)
	createList := make([]*dto.CreateTopicDto, 0, len(addFiles))
	dtoUpdateList := make([]*dto.CreateTopicDto, 0, len(addFiles))

	for _, file := range addFiles {
		file.Status = 1
		createTopicDto := UnsConverter.Po2Dto(file)

		var dbF *dao.UnsNamespace
		if temp, exists := dbFiles[file.ID]; exists {
			dbF = temp
		} else {
			dbF = existsUns[file.Alias]
		}

		if labels, exists := unsPoLabels[file.ID]; exists {
			labels.SetDto(createTopicDto)
		}
		if dbF != nil && dbF.Status == OK {
			switch file.PathType {
			case constants.PathTypeFile:
				createTopicDto.FieldsChanged = !base.EqualsF(file.Fields, dbF.Fields, func(a, b *dto.FieldDefine) bool {
					return a.Name == b.Name && a.Type == b.Type
				})
			case constants.PathTypeDir:
				/*	if file.LayRec != dbF.LayRec {
					createTopicDto.OldPath = dbF.Path
					createTopicDto.OldLayRec = dbF.LayRec
					protocolStr := file.Protocol
					if protocolStr == "" || protocolStr[0] != '{' {
						protocolStr = "{}"
						file.Protocol = protocolStr
					}
					var protocol map[string]interface{}
					if err := json.Unmarshal([]byte(protocolStr), &protocol); err == nil {
						protocol["oldLay"] = dbF.LayRec
						if updatedProtocol, err := json.Marshal(protocol); err == nil {
							file.Protocol = string(updatedProtocol)
						}
					}
				}*/
			}
			dtoUpdateList = append(dtoUpdateList, createTopicDto)
		} else {
			if dbF != nil {
				// 逻辑删除后重建，设置默认值用来update
				if file.MountType == nil {
					zero := int16(0)
					file.MountType = &zero
				}
				if file.WithFlags == nil {
					zero := int32(0)
					file.WithFlags = &zero
				}
				if file.ExtendFieldFlags == nil {
					zero := int32(0)
					file.ExtendFieldFlags = &zero
				}
				if file.RefUns == nil {
					file.RefUns = make(dao.RefUns)
				}
			}
			file.CreateAt = createTime
			file.UpdateAt = createTime
			createTopicDto.CreateAt = createTime
			createTopicDto.UpdateAt = createTime
			createList = append(createList, createTopicDto)
		}
	}

	for _, po := range rs.updateList {
		id := po.ID
		po.Status = 1
		if _, exists := addFiles[id]; !exists {
			dtoUpdateList = append(dtoUpdateList, UnsConverter.Po2Dto(po))
		}
	}

	/*	if refUpdates != nil && len(refUpdates) > 0 {
		for _, refPo := range refUpdates {
			id := refPo.ID
			if po, exists := dbFiles[id]; exists {
				po.Status = 1
				rs.UpdateList = append(rs.UpdateList, po)
				if _, exists := addFiles[id]; !exists {
					dtoUpdateList = append(dtoUpdateList, UnsConverter.Po2Dto(po, false))
				}
			}
		}
	}*/
	var unsLabels = base.MapMapValues(unsPoLabels, func(upl *bo.UnsPoLabels) bo.UnsLabels {
		return upl
	})
	err = u.saveBatchAndSendEvent(ctx, createTime, &args, rs.insertList, rs.updateList,
		createList, dtoUpdateList, unsLabels)
	if err != nil {
		errTipMap["0"] = err.Error()
	}
	return errTipMap
}
func (u *UnsAddService) saveBatchAndSendEvent(
	ctx context.Context,
	createTime time.Time,
	args *bo.CreateModelInstancesArgs,
	insertList []*dao.UnsNamespace,
	updateList []*dao.UnsNamespace,
	notifyCreateList []*dto.CreateTopicDto,
	notifyUpdateList []*dto.CreateTopicDto,
	unsLabels []bo.UnsLabels) error {

	dataSrcFiles := base.GroupBy(notifyCreateList, func(e *dto.CreateTopicDto) common.SrcJdbcType {
		return e.DataSrcID
	})
	tx := dao.GetDb(ctx).Begin()
	ctx = dao.SetDb(ctx, tx)
	defer func() {
		if r := recover(); r != nil {
			u.log.Error("SaveUnsPanic:", r)
			tx.Rollback()
		}
	}()
	labelPos, err := u.unsLabelService.MakeUnsLabels(ctx, unsLabels, createTime)
	if err == nil {
		if len(insertList) > 0 {
			err = u.unsMapper.MultiInsert(tx, insertList)
			u.log.Debug("insertUns:", len(insertList), err)
		}
		if err == nil && len(updateList) > 0 {
			err = u.unsMapper.MultiUpdate(tx, updateList)
			u.log.Debug("updateUns:", len(insertList), err)
		}
	}
	if err == nil {
		err = spring.PublishEvent(&event.BatchCreateTableEvent{
			ApplicationEvent: event.ApplicationEvent{Context: ctx},
			FlowName:         args.FlowName,
			FromImport:       args.FromImport,
			Topics:           dataSrcFiles,
			Updates:          notifyUpdateList,
			Labels:           base.Map(labelPos, UnsConverter.Label2Uns),
			DelegateAware:    getEventStatusCallback(args.StatusConsumer),
		})
	}
	if err != nil {
		u.log.Error("SaveUnsErr:", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}
func (u *UnsAddService) CreateModelInstance(ctx context.Context, topicDto *dto.CreateTopicDto) *types.StringResult {
	result := &types.StringResult{BaseResult: types.BaseResult{Code: 200, Msg: "ok"}}
	db := dao.GetDb(ctx)
	// 处理父文件夹ID
	if topicDto.ParentID != nil && *topicDto.ParentID != 0 && topicDto.ParentAlias == "" {
		folder, err := u.unsMapper.SelectById(db, *topicDto.ParentID)
		if err != nil || folder == nil {
			result.Code = 400
			result.Msg = I18nUtils.GetMessage("uns.folder.not.found")
			return result
		}

		//if folder.MountType != nil && MountSourceType.IsCollectorMountSource(*folder.MountType) {
		//	return &JsonResult[string]{Code: 400, Message: I18nUtils.GetMessage("uns.mount.folder.operate")}
		//}
		topicDto.ParentAlias = folder.Alias
	}

	// TODO 是文件夹并且需要创建模板
	if topicDto.PathType == constants.PathTypeDir && topicDto.CreateTemplate != nil && *topicDto.CreateTemplate {
		//templateVo := &CreateTemplateVo{
		//	Name:   topicDto.Name,
		//	Fields: topicDto.Fields,
		//}
		//templateResult := u.unsTemplateService.CreateTemplate(templateVo)
		//if templateResult.Code != 200 {
		//	result.Code = 400
		//	result.Message = templateResult.Message
		//	return result
		//}
		//modelId, _ := strconv.ParseInt(templateResult.Data, 10, 64)
		//topicDto.ModelId = &modelId
	}

	topicDto.Index = 0

	args := bo.CreateModelInstancesArgs{
		Topics:              []*dto.CreateTopicDto{topicDto},
		FromImport:          false,
		ThrowModelExistsErr: true,
	}

	// 设置计算类型不添加流程
	if topicDto.DataType != nil && (*topicDto.DataType == constants.CalculationHistType || *topicDto.DataType == constants.CalculationRealType) {
		falseVal := false
		topicDto.AddFlow = &falseVal
	}
	rs := u.CreateModelAndInstancesInner(ctx, args)
	if rs != nil && len(rs) > 0 {
		// 将map中的错误信息拼接成字符串
		var errorMessages []string
		for _, msg := range rs {
			errorMessages = append(errorMessages, msg)
		}
		result.Code = 400
		result.Msg = strings.Join(errorMessages, ", ")
	} else {
		if topicDto.ID > 0 {
			result.Data = strconv.FormatInt(topicDto.ID, 10)
		}
	}

	return result
}
func (u *UnsAddService) CreateModelAndInstance(ctx context.Context, topicDtos []*dto.CreateTopicDto, fromImport bool) map[string]string {
	taskID := fmt.Sprintf("%p_%d", &topicDtos, len(topicDtos)) // 使用指针地址作为任务ID

	args := bo.CreateModelInstancesArgs{
		Topics:              topicDtos,
		FromImport:          fromImport,
		ThrowModelExistsErr: false,
		FlowName:            time.Now().Format("20060102150405"), // yyyyMMddHHmmss
	}
	// 设置状态消费者
	args.StatusConsumer = func(runningStatus *common.RunningStatus) {
		if runningStatus.SpendMills != nil {
			i := runningStatus.I
			n := runningStatus.N
			task := runningStatus.Task
			u.log.Infof("[%u] %d/%d 已处理， %u：耗时%d ms", taskID, i, n, task, *runningStatus.SpendMills)
		}
	}

	parentAliasMap := make(map[string]string)

	// 处理每个topicDto
	for i, topicDto := range topicDtos {
		topicDto.Batch = 0
		topicDto.Index = i

		if topicDto.ParentAlias != "" {
			batchIndex := topicDto.GainBatchIndex()
			parentAliasMap[batchIndex] = topicDto.ParentAlias
		}
	}
	//db := dao.GetDb(ctx)
	// 检查挂载文件夹
	/*	if len(parentAliasMap) > 0 {
		// 收集所有父别名
		var parentAliases = base.MapValues(parentAliasMap)
		parentUnsList, err := u.unsMapper.ListByAlias(db, parentAliases)
		if err == nil && len(parentUnsList) > 0 {
			mountAlias := make(map[string]bool)

			for _, uns := range parentUnsList {
				if uns.MountType != nil && MountSourceType.IsCollectorMountSource(*uns.MountType) {
					mountAlias[uns.Alias] = true
				}
			}

			if len(mountAlias) > 0 {
				rs := make(map[string]string)
				for batchIndex, alias := range parentAliasMap {
					if mountAlias[alias] {
						rs[batchIndex] = I18nUtils.GetMessage("uns.mount.folder.operate")
					}
				}
				if len(rs) > 0 {
					return rs
				}
			}
		}
	}*/

	// 处理标签
	labelsMap := make(map[string][]string)
	for _, topicDto := range topicDtos {
		if topicDto.LabelNames != nil && len(topicDto.LabelNames) > 0 {
			labelsMap[topicDto.Alias] = topicDto.LabelNames
		}
	}
	args.LabelsMap = labelsMap

	rs := u.CreateModelAndInstancesInner(ctx, args)
	u.log.Infof("[%u] UNS 处理完毕.", taskID)
	return rs
}
