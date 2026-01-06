package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	sysconfig "backend/internal/common/config"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/label/service"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"encoding/json"
	"errors"
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
	sysConfig       *sysconfig.SystemConfig
	unsLabelService *service.UnsLabelService
	removeService   *UnsRemoveService
	unsCalcService  UnsCalcService
}

func init() {
	spring.RegisterLazy[*UnsAddService](func() *UnsAddService {
		return &UnsAddService{
			log:             logx.WithContext(context.Background()),
			sysConfig:       spring.GetBean[*sysconfig.SystemConfig](),
			unsLabelService: spring.GetBean[*service.UnsLabelService](),
			removeService:   spring.GetBean[*UnsRemoveService](),
		}
	})
}
func (u *UnsAddService) CreateModelAndInstancesInner(ctx context.Context, args bo.CreateModelInstancesArgs) (errTipMap map[string]string) {
	{
		db := dao.GetDb(ctx)
		ctx = dao.SetDb(ctx, db)
	}
	// 对文件进行归类
	errTipMap = make(map[string]string, len(args.Topics))
	for i, topic := range args.Topics {
		if topic.Index == 0 {
			topic.Index = i
		}
	}
	if u.sysConfig.EnableAutoCategorization {
		args.Topics = u.appendCategoryFolders(ctx, args.Topics, errTipMap)
	}
	u.log.Debugf("[%d] args: %+v", len(args.Topics), args)
	topicDtos := args.Topics
	pathMap := initParamsUns(topicDtos, errTipMap)
	if len(pathMap) == 0 {
		u.log.Info("不存在任何文件夹或文件, 无法继续保存")
		return errTipMap
	}
	dbFiles := make(map[int64]*dao.UnsNamespace)
	var existsUns = make(map[string]*dao.UnsNamespace, int(float64(len(args.Topics))*1.3)+5)
	{
		ids := make(map[int64]bool)
		aliasSet := make(map[string]bool)
		for _, vs := range pathMap {
			addAlias(base.MapValues(vs), aliasSet, ids)
		}
		var err error
		err = u.listUnsByAliasAndIds(ctx, base.MapKeys(aliasSet), base.MapKeys(ids), dbFiles, existsUns)
		if err != nil {
			u.log.Error("QueryUnsError:", err)
			errTipMap["0"] = "QueryUnsError:" + err.Error()
			return errTipMap
		}
		if len(existsUns) > 0 {
			newIds := make(map[int64]bool)
			for _, uns := range existsUns {
				pid := uns.ParentId
				if pid != nil && !ids[*pid] {
					newIds[*pid] = true
				}
				mid := uns.ModelId
				if mid != nil && !ids[*mid] {
					newIds[*mid] = true
				}
			}
			if len(newIds) > 0 {
				db := dao.GetDb(ctx)
				for _, idList := range base.Partition(base.MapKeys(newIds), constants.SQLBatchSize) {
					unsPos, err := u.unsMapper.SelectByIds(db, idList)
					if err != nil {
						u.log.Error("QueryUnsError2:", err)
						errTipMap["0"] = "QueryUnsError2:" + err.Error()
						return errTipMap
					}
					addDbPo(unsPos, dbFiles, existsUns)
				}
			}
		}
	}
	for _, vs := range pathMap {
		tryFillIdOrAlias(vs, existsUns, dbFiles, errTipMap)
	}
	paramFiles, paramFolders := pathMap[constants.PathTypeFile], pathMap[constants.PathTypeDir]
	addFiles := make(map[int64]*dao.UnsNamespace)
	aliasMap := make(map[string]*dao.UnsNamespace)
	folders := base.MapValues(paramFolders)
	if len(folders) > 1 {
		var loopFolders []*types.CreateTopicDto
		folders, loopFolders = base.SorByDependency(folders, func(a, b *types.CreateTopicDto) bool {
			return a.Index < b.Index
		}, func(t *types.CreateTopicDto) string {
			return t.Alias
		}, func(t *types.CreateTopicDto) string {
			return base.P2v(t.ParentAlias)
		})
		if len(loopFolders) > 0 {
			msg := I18nUtils.GetMessage("uns.circularDependency")
			for _, folder := range loopFolders {
				errTipMap[folder.GainBatchIndex()] = msg + ": " + folder.Alias
			}
		}
	}
	createTime := time.Now()
	allUns := func(alias string) *dao.UnsNamespace {
		unsPo := aliasMap[alias]
		if unsPo == nil {
			unsPo = existsUns[alias]
		}
		return unsPo
	}
	unsPoLabels := make(map[int64]*bo.UnsPoLabels, len(paramFiles)+len(paramFolders))
	var deleteFiles []*dao.UnsNamespace

	indirectUpdates := make([]*dao.UnsNamespace, 0)
	addUpdate := func(uns *dao.UnsNamespace) {
		indirectUpdates = append(indirectUpdates, uns)
	}
	u.itrFiles(ctx, addUpdate, args.SkipWhenExists, base.MapValues(pathMap[constants.PathTypeTemplate]), createTime, allUns, dbFiles, deleteFiles, errTipMap, addFiles, aliasMap, unsPoLabels)
	u.itrFiles(ctx, addUpdate, args.SkipWhenExists, folders, createTime, allUns, dbFiles, deleteFiles, errTipMap, addFiles, aliasMap, unsPoLabels)
	files := base.MapValues(pathMap[constants.PathTypeFile])
	sort.Sort(indexedFiles(files))
	u.itrFiles(ctx, addUpdate, args.SkipWhenExists, files, createTime, allUns, dbFiles, deleteFiles, errTipMap, addFiles, aliasMap, unsPoLabels)

	//TODO 计算，引用，聚合等类型的 校验和处理
	aliasToId(addFiles, allUns, pathMap)
	// 找出 parentAlias 或 name 有修改的最高层目录，后面需要获取它的整个子树，为更新 layRec 做准备
	err := u.tryAddLayRecOrPathChangedChildren(ctx, addFiles, dbFiles, existsUns)
	if err != nil {
		errTipMap["0"] = err.Error()
		return errTipMap
	}
	validDBFiles := base.MapFilter(dbFiles, func(dbPo *dao.UnsNamespace) bool {
		return base.P2v(dbPo.Status) == OK
	})
	rs := setLayRecAndPath(createTime, addFiles, validDBFiles)
	createList := make([]*types.CreateTopicDto, 0, len(addFiles))
	dtoUpdateList := make([]*types.CreateTopicDto, 0, len(addFiles))

	for _, file := range addFiles {
		file.Status = &OK
		createTopicDto := UnsConverter.Po2Dto(file)

		var dbF *dao.UnsNamespace
		if temp, exists := dbFiles[file.Id]; exists {
			dbF = temp
		} else {
			dbF = existsUns[file.Alias]
		}

		if labels, exists := unsPoLabels[file.Id]; exists {
			labels.SetDto(createTopicDto)
		}
		if dbF != nil && base.P2v(dbF.Status) == OK {
			dtoUpdateList = append(dtoUpdateList, createTopicDto)
		} else {
			if dbF != nil {
				// 逻辑删除后重建，设置默认值用来update
				if file.MountType == nil {
					zero := int16(0)
					file.MountType = &zero
					createTopicDto.MountType = file.MountType
				}
				if file.WithFlags == nil {
					zero := int32(0)
					file.WithFlags = &zero
					createTopicDto.WithFlags = file.WithFlags
				}
				if file.ExtendFieldFlags == nil {
					zero := int32(0)
					file.ExtendFieldFlags = &zero
				}
				if file.RefUns == nil {
					file.RefUns = make(dao.RefUns)
				}
				id := file.Id
				rs.updateList[id] = file
				delete(rs.insertList, id)
			}
			file.CreateAt = createTime
			file.UpdateAt = createTime
			createTopicDto.CreateAt = createTime.UnixMilli()
			createTopicDto.UpdateAt = createTopicDto.CreateAt
			createList = append(createList, createTopicDto)
		}
	}

	for _, po := range rs.updateList {
		id := po.Id
		po.Status = &OK
		if _, exists := addFiles[id]; !exists {
			dtoUpdateList = append(dtoUpdateList, UnsConverter.Po2Dto(po))
		}
	}
	if len(indirectUpdates) > 0 {
		for _, update := range indirectUpdates {
			rs.updateList[update.Id] = update
			dtoUpdateList = append(dtoUpdateList, UnsConverter.Po2Dto(update))
		}
	}

	/*	if refUpdates != nil && len(refUpdates) > 0 {
		for _, refPo := range refUpdates {
			id := refPo.Id
			if po, exists := dbFiles[id]; exists {
				po.Status = 1
				rs.UpdateList = append(rs.UpdateList, po)
				if _, exists := addFiles[id]; !exists {
					dtoUpdateList = append(dtoUpdateList, UnsConverter.Po2Dto(po, false))
				}
			}
		}
	}*/
	u.log.Infof("addFiles:%d,db:%d, createList.size=%d, updateList.size=%d\n", len(addFiles), len(dbFiles), len(rs.insertList), len(rs.updateList))
	var unsLabels = base.MapMapValues(unsPoLabels, func(upl *bo.UnsPoLabels) bo.UnsLabels {
		return upl
	})
	err = u.saveBatchAndSendEvent(ctx, createTime, &args, base.MapValues(rs.insertList), base.MapValues(rs.updateList),
		createList, dtoUpdateList, deleteFiles, unsLabels)
	if err != nil {
		errTipMap["0"] = err.Error()
	}
	return errTipMap
}

type indexedFiles []*types.CreateTopicDto

func (x indexedFiles) Len() int           { return len(x) }
func (x indexedFiles) Less(i, j int) bool { return x[i].Index < x[j].Index }
func (x indexedFiles) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (u *UnsAddService) itrFiles(
	ctx context.Context,
	addUpdate func(*dao.UnsNamespace),
	skipWhenExists bool,
	vs []*types.CreateTopicDto,
	createTime time.Time, allUns func(alias string) *dao.UnsNamespace,
	dbFiles map[int64]*dao.UnsNamespace, deleteFiles []*dao.UnsNamespace,
	errTipMap map[string]string,
	addFiles map[int64]*dao.UnsNamespace,
	aliasMap map[string]*dao.UnsNamespace,
	unsPoLabels map[int64]*bo.UnsPoLabels) {
	if len(vs) == 0 {
		return
	}
	for _, DTO := range vs {
		po, exists := u.trySetId(ctx, skipWhenExists, createTime, DTO, allUns, dbFiles, addUpdate, &deleteFiles, errTipMap)
		if po != nil {
			if exists && skipWhenExists {
				aliasMap[po.Alias] = po
				continue
			}
			addFiles[po.Id] = po
			aliasMap[po.Alias] = po
			if DTO.LabelNames != nil {
				_, exists := dbFiles[po.Id]
				unsPoLabels[po.Id] = bo.NewUnsPoLabels(po, exists, DTO.LabelNames)
			}
		}
	}
}
func (u *UnsAddService) saveBatchAndSendEvent(
	ctx context.Context,
	createTime time.Time,
	args *bo.CreateModelInstancesArgs,
	insertList []*dao.UnsNamespace,
	updateList []*dao.UnsNamespace,
	notifyCreateList []*types.CreateTopicDto,
	notifyUpdateList []*types.CreateTopicDto,
	deleteFiles []*dao.UnsNamespace,
	unsLabels []bo.UnsLabels) error {

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
		if len(notifyCreateList)+len(notifyUpdateList)+len(labelPos) > 0 {
			createFiles := base.GroupBy(notifyCreateList, pathTypeGroupBy)
			notifyUpdate := base.GroupBy(notifyUpdateList, pathTypeGroupBy)
			if len(labelPos) > 0 {
				notifyUpdate[constants.PathTypeLabel] = base.Map(labelPos, UnsConverter.Label2Uns)
			}
			err = spring.PublishEvent(&event.BatchCreateTableEvent{
				ApplicationEvent: event.ApplicationEvent{Context: ctx},
				FlowName:         args.FlowName,
				FromImport:       args.FromImport,
				Creates:          createFiles,
				Updates:          notifyUpdate,
				DelegateAware:    getEventStatusCallback(args.StatusConsumer),
			})
		}
		if len(deleteFiles) > 0 {
			delRs, delEr := u.removeService.Remove(ctx, types.RemoveUnsOptions{
				WithFlow:    base.V2p(args.FlowName != ""),
				RemoveRefer: &TRUE,
			}, deleteFiles)
			if delEr != nil {
				err = delEr
			} else if delRs != nil && delRs.Code != 200 && delRs.Msg != "" {
				err = errors.New(delRs.Msg)
			}
		}
	}
	if err != nil {
		u.log.Error("SaveUnsErr:", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}
func pathTypeGroupBy(e *types.CreateTopicDto) int16 {
	return e.PathType
}
func (u *UnsAddService) CreateModelInstance(ctx context.Context, topicDto *types.CreateTopicDto) *types.CreateUnsResp {
	result := &types.CreateUnsResp{BaseResult: types.BaseResult{Code: 200, Msg: "ok"}}
	db := dao.GetDb(ctx)
	// 处理父文件夹ID
	if topicDto.ParentId != nil && *topicDto.ParentId != 0 && topicDto.ParentAlias == nil {
		folder, err := u.unsMapper.SelectById(db, *topicDto.ParentId)
		if err != nil {
			result.Code = 500
			result.Msg = err.Error()
			return result
		} else if folder == nil {
			result.Code = 400
			result.Msg = I18nUtils.GetMessage("uns.folder.not.found") + ":id=" + strconv.Itoa(int(*topicDto.ParentId))
			return result
		}

		//if folder.MountType != nil && MountSourceType.IsCollectorMountSource(*folder.MountType) {
		//	return &JsonResult[string]{Code: 400, Message: I18nUtils.GetMessage("uns.mount.folder.operate")}
		//}
		topicDto.ParentAlias = &folder.Alias
	}
	unsList := make([]*types.CreateTopicDto, 0, 2)
	// 是文件夹 并且需要创建模板
	if topicDto.PathType == constants.PathTypeDir && base.P2v(topicDto.CreateTemplate) {
		var templateAlias = "T" + PathUtil.GenerateAliasWithRandom(topicDto.Name)
		templateVo := &types.CreateTopicDto{
			PathType: constants.PathTypeTemplate,
			Alias:    templateAlias,
			Name:     topicDto.Name,
			Fields:   topicDto.Fields,
		}
		unsList = append(unsList, templateVo)
		topicDto.ModelAlias = &templateVo.Alias
	}

	topicDto.Index = 0
	unsList = append(unsList, topicDto)

	args := bo.CreateModelInstancesArgs{
		Topics:              unsList,
		FromImport:          false,
		ThrowModelExistsErr: true,
		StatusConsumer: func(status *common.RunningStatus) {
			bs, _ := json.Marshal(status)
			u.log.Infof("创建UNS[%s - %s], 进度：%v", topicDto.Name, topicDto.Alias, string(bs))
		},
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
	} else if topicDto.Id > 0 {
		result.Data = types.CreateUnsResult{Id: strconv.FormatInt(topicDto.Id, 10)}
		if topicDto.ParentId != nil {
			result.Data.ParentId = strconv.FormatInt(*topicDto.ParentId, 10)
		}
	}

	return result
}
func (u *UnsAddService) CreateModelAndInstance(ctx context.Context, topicDtos []*types.CreateTopicDto, fromImport bool) map[string]string {
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

	//db := dao.GetDb(ctx)
	// 检查挂载文件夹
	/*
			parentAliasMap := make(map[string]string)

			// 处理每个topicDto
			for i, topicDto := range topicDtos {
				topicDto.Batch = 0
				topicDto.Index = i

				if parentAlias := topicDto.ParentAlias; parentAlias != nil {
					batchIndex := topicDto.GainBatchIndex()
					parentAliasMap[batchIndex] = *parentAlias
				}
			}

		if len(parentAliasMap) > 0 {
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
	rs := u.CreateModelAndInstancesInner(ctx, args)
	u.log.Infof("[%u] UNS 处理完毕.", taskID)
	return rs
}
