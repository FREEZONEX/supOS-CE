package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/event"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
)

type UnsLabelService struct {
	unsMapper      dao.UnsNamespaceRepo
	labelMapper    dao.UnsLabelRepo
	labelRefMapper dao.UnsLabelRefRepo
}

func init() {
	spring.RegisterBean[*UnsLabelService](&UnsLabelService{})
}
func (l *UnsLabelService) AllLabel(ctx context.Context, req *types.UnsLabelListReq) (resp *types.UnsLabelListResp, err error) {
	dbFilter := dao.UnsLabelFilter{
		LabelName: req.Key,
	}
	db := dao.GetDb(ctx)
	list, err := l.labelMapper.FindByFilter(db, dbFilter, nil)
	if err != nil {
		return nil, err
	}
	resp = &types.UnsLabelListResp{
		List: make([]*types.UnsLabel, 0),
	}
	for _, item := range list {
		resp.List = append(resp.List, &types.UnsLabel{
			ID:        item.ID,
			LabelName: item.LabelName,
			CreateAt:  item.CreateAt.UnixMilli(),
		})
	}
	return
}

func (l *UnsLabelService) Create(ctx context.Context, req *types.UnsLabelCreateReq) (resp *types.WithID, err error) {
	// 参数校验
	if strings.TrimSpace(req.LabelName) == "" {
		return nil, errors.Parameter.WithMsg("labelName不能为空")
	}

	// 写入数据库
	data := &dao.UnsLabel{
		LabelName: req.LabelName,
	}
	db := dao.GetDb(ctx)
	if err = l.labelMapper.Insert(db, data); err != nil {
		return nil, err
	}

	// 返回新ID
	return &types.WithID{ID: data.ID}, nil
}
func (l *UnsLabelService) Delete(ctx context.Context, req *types.WithID) (resp *types.Empty, err error) {
	id := req.ID
	if id <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	err = dao.GetDb(ctx).Transaction(func(tx *gorm.DB) (er error) {
		if err = l.labelMapper.Delete(tx, id); err != nil {
			return err
		}
		unsIds, er := l.labelRefMapper.ListUnsIds(tx, id)
		if len(unsIds) > 0 {
			updateTime := time.Now()
			er = l.labelMapper.DeleteRefByLabelId(tx, id)
			if er != nil {
				return er
			}
			for _, parUnsIds := range base.Partition[int64](unsIds, 500) {
				_, er = l.unsMapper.UnlinkLabelsByIds(tx, id, parUnsIds, updateTime)
				if er != nil {
					return er
				}
			}
		}
		return er
	})

	return &types.Empty{}, err
}
func (l *UnsLabelService) Detail(ctx context.Context, req *types.WithID) (resp *types.UnsLabel, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	db := dao.GetDb(ctx)
	item, err := l.labelMapper.FindOne(db, req.ID)
	if err != nil {
		return nil, err
	}
	return &types.UnsLabel{
		ID:        item.ID,
		LabelName: item.LabelName,
		CreateAt:  item.CreateAt.UnixMilli(),
	}, nil
}
func (l *UnsLabelService) Update(ctx context.Context, req *types.UnsLabel) (resp *types.Empty, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	if strings.TrimSpace(req.LabelName) == "" {
		return nil, errors.Parameter.WithMsg("labelName不能为空")
	}
	db := dao.GetDb(ctx)
	item, err := l.labelMapper.FindOne(db, req.ID)
	if err != nil {
		return nil, err
	}
	item.LabelName = req.LabelName
	if err = l.labelMapper.Update(db, item); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (l *UnsLabelService) CreateBatch(ctx context.Context, labels []string) (err error) {
	db := dao.GetDb(ctx)
	labelsPos, er := l.labelMapper.FindByNames(db, labels)
	if er != nil {
		return er
	}
	labelMap := base.MapArrayToMap[*dao.UnsLabel, string, *dao.UnsLabel](labelsPos, func(e *dao.UnsLabel) (ok bool, k string, v *dao.UnsLabel) {
		return true, e.LabelName, e
	})
	insertList := make([]*dao.UnsLabel, 0, len(labels))
	updateTime := time.Now()
	for _, label := range labels {
		po, has := labelMap[label]
		if !has {
			po = &dao.UnsLabel{ID: common.NextId(), LabelName: label}
			insertList = append(insertList, po)
			po.CreateAt = updateTime
		}
	}
	if len(insertList) > 0 {
		err = l.labelMapper.MultiInsert(db, insertList)
	}
	return err
}
func (l *UnsLabelService) MakeUnsLabels(ctx context.Context, unsLabels []bo.UnsLabels, createTime time.Time) (rs []*dao.UnsLabel, er error) {
	if len(unsLabels) == 0 {
		return nil, nil
	}
	var resetUnsIds []int64
	labelUnsMap := make(map[string][]bo.UnsLabels)
	for _, unsLabel := range unsLabels {
		if unsLabel.IsResetLabels() {
			if resetUnsIds == nil {
				resetUnsIds = make([]int64, 0, len(unsLabels))
			}
			resetUnsIds = append(resetUnsIds, unsLabel.UnsId())
		}
		labelNames := unsLabel.LabelNames()
		if len(labelNames) > 0 {
			for _, label := range labelNames {
				labelUnsMap[label] = append(labelUnsMap[label], unsLabel)
			}
		}
	}
	db := dao.GetDb(ctx)
	if len(resetUnsIds) > 0 {
		er := l.labelRefMapper.DeleteByUnsIds(db, resetUnsIds)
		if er != nil {
			return nil, er
		}
	}
	labels := base.MapKeys[string](labelUnsMap)
	allLabels := make([]*dao.UnsLabel, 0, len(labels))
	var saveLabels []*dao.UnsLabel
	saveLabelRef := make([]*dao.UnsLabelRef, 0, len(labels))
	if len(labels) > 0 {
		existLabels, er := l.labelMapper.FindByNames(db, labels)
		if er != nil {
			return nil, er
		}
		existLabelMap := base.MapArrayToMap[*dao.UnsLabel, string, *dao.UnsLabel](existLabels, func(e *dao.UnsLabel) (ok bool, k string, v *dao.UnsLabel) {
			return true, e.LabelName, e
		})
		for _, label := range labels {
			existLabel := existLabelMap[label]
			labelPo := existLabel
			if labelPo == nil {
				labelPo = &dao.UnsLabel{LabelName: label, CreateAt: createTime}
			}
			if existLabel == nil {
				// 新增标签
				labelPo.ID = common.NextId()
				if saveLabels == nil {
					saveLabels = make([]*dao.UnsLabel, 0, len(labels))
				}
				saveLabels = append(saveLabels, labelPo)
				allLabels = append(allLabels, labelPo)
			}
			labelId := labelPo.ID
			unsLabelsList := labelUnsMap[label]
			for _, ul := range unsLabelsList {
				ul.SetLabelId(label, labelId)
			}
			saveLabelRef = append(saveLabelRef, base.Map[bo.UnsLabels, *dao.UnsLabelRef](unsLabelsList, func(e bo.UnsLabels) *dao.UnsLabelRef {
				return &dao.UnsLabelRef{LabelID: labelId, UnsID: e.UnsId()}
			})...)
		}
	}
	if len(saveLabels) > 0 {
		er = l.labelMapper.MultiInsert(db, saveLabels)
		if er != nil {
			return nil, er
		}
	}
	if len(saveLabelRef) > 0 {
		er = l.labelRefMapper.SaveOrIgnore(db, saveLabelRef)
	}
	return allLabels, er
}
func (s *UnsLabelService) MakeLabels(ctx context.Context, unsId int64, labelList []*types.LabelVo) error {
	db := dao.GetDb(ctx)
	uns, er := s.unsMapper.SelectById(db, unsId)
	if uns == nil {
		return nil
	} else if er != nil {
		return er
	}

	labelIdNameMap := make(map[int64]string)
	uns.LabelIds = labelIdNameMap
	er = s.labelMapper.DeleteRefByUnsId(db, unsId)
	if er != nil {
		return er
	}

	if len(labelList) > 0 {
		var noNames = base.MapArrayToMap[*types.LabelVo, int64, *types.LabelVo](labelList, func(e *types.LabelVo) (ok bool, k int64, v *types.LabelVo) {
			if v.ID != 0 && v.LabelName == "" {
				ok, k, v = true, v.ID, e
			}
			return
		})

		if len(noNames) > 0 {
			ids := base.MapKeys[int64](noNames)
			labels, er := s.labelMapper.ListByIds(db, ids)
			if er != nil {
				return er
			} else if len(labels) > 0 {
				for _, lb := range labels {
					vo := noNames[lb.ID]
					if vo != nil {
						vo.LabelName = lb.LabelName
					}
				}
			}
		}
		labels := make([]*dao.UnsLabel, 0, len(labelList))
		refs := make([]*dao.UnsLabelRef, 0, len(labelList))
		now := time.Now()
		for _, labelVo := range labelList {
			lid := labelVo.ID
			var ref *dao.UnsLabelRef
			if lid != 0 {
				ref = &dao.UnsLabelRef{LabelID: lid, UnsID: unsId}
			} else {
				// 创建标签
				// 假设创建成功并返回 Id
				label := &dao.UnsLabel{ID: common.NextId(), LabelName: labelVo.LabelName, CreateAt: now}
				ref = &dao.UnsLabelRef{LabelID: label.ID, UnsID: unsId}
				labels = append(labels, label)
			}
			labelIdNameMap[ref.LabelID] = labelVo.LabelName
			refs = append(refs, ref)
		}
		if len(labels) > 0 {
			er = s.labelMapper.MultiInsert(db, labels)
			if er != nil {
				return er
			}
		}
		if len(refs) > 0 {
			er = s.labelRefMapper.MultiInsert(db, refs)
			if er != nil {
				return er
			}
		}
	}
	uns.UpdateAt = time.Now()
	er = s.unsMapper.Update(db, uns)
	return er
}
func (s *UnsLabelService) CancelLabel(ctx context.Context, unsId int64, labelIds []int64) error {
	db := dao.GetDb(ctx)
	uns, er := s.unsMapper.SelectById(db, unsId)
	if uns == nil {
		return nil
	} else if er != nil {
		return er
	}
	er = s.labelRefMapper.DeleteByUnsIdAndLabelIds(db, uns.Id, labelIds)
	if er != nil {
		return er
	}
	leftLabels, er := s.labelMapper.ListByUnsId(db, unsId)
	if er != nil {
		return er
	}
	labelIdMap := make(map[int64]string)
	for _, label := range leftLabels {
		labelIdMap[label.ID] = label.LabelName
	}
	now := time.Now()
	updatePo := &dao.UnsNamespace{}
	updatePo.LabelIds = labelIdMap
	updatePo.Id = unsId
	updatePo.UpdateAt = now
	er = s.unsMapper.Update(db, updatePo)
	if er != nil {
		return er
	}
	return nil
}
func (s *UnsLabelService) CancelLabelByNames(ctx context.Context, unsAlias string, labelNames []string) error {
	db := dao.GetDb(ctx)
	uns, er := s.unsMapper.GetByAlias(db, unsAlias)
	if uns == nil {
		return errors.NewCodeError(400, I18nUtils.GetMessage("uns.file.not.exist"))
	} else if er != nil {
		return er
	}
	labelList, er := s.labelMapper.FindByNames(db, labelNames)
	if len(labelList) == 0 {
		return nil
	} else if er != nil {
		return nil
	}
	labelIds := base.Map[*dao.UnsLabel, int64](labelList, func(e *dao.UnsLabel) int64 {
		return e.ID
	})
	ctx = dao.SetDb(ctx, db)
	return s.CancelLabel(ctx, uns.Id, labelIds)
}

// OnEventRemoveTopicsEvent 处理UNS 删除事件
func (s *UnsLabelService) OnEventRemoveTopicsEvent(event event.RemoveTopicsEvent) (er error) {
	labelIds := base.FilterAndFlatMap[*types.CreateTopicDto, int64](event.Topics, func(e *types.CreateTopicDto) (vs []int64, ok bool) {
		if len(e.LabelIDs) > 0 {
			vs, ok = base.MapKeys(e.LabelIDs), true
		}
		return
	})
	if len(labelIds) == 0 {
		return nil
	}
	db := dao.GetDb(event.Context)
	for _, partLabelIds := range base.Partition(labelIds, 1000) {
		er = s.labelRefMapper.DeleteByLabelIds(db, partLabelIds)
		if er != nil {
			return er
		}
	}
	return
}
