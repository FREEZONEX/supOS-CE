package service

import (
	"backend/internal/common"
	"backend/internal/logic/supos/uns/uns/bo"
	"backend/internal/repo/relationDB"
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
	unsMapper      relationDB.UnsNamespaceRepo
	labelMapper    relationDB.UnsLabelRepo
	labelRefMapper relationDB.UnsLabelRefRepo
}

func init() {
	spring.RegisterBean[*UnsLabelService](&UnsLabelService{})
}
func (l *UnsLabelService) AllLabel(ctx context.Context, req *types.UnsLabelListReq) (resp *types.UnsLabelListResp, err error) {
	dbFilter := relationDB.UnsLabelFilter{
		LabelName: req.Key,
	}
	db := relationDB.GetDb(ctx)
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
	data := &relationDB.UnsLabel{
		LabelName: req.LabelName,
	}
	db := relationDB.GetDb(ctx)
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
	err = relationDB.GetDb(ctx).Transaction(func(tx *gorm.DB) (er error) {
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
	db := relationDB.GetDb(ctx)
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
	db := relationDB.GetDb(ctx)
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
	db := relationDB.GetDb(ctx)
	labelsPos, er := l.labelMapper.FindByNames(db, labels)
	if er != nil {
		return er
	}
	labelMap := base.MapArrayToMap[*relationDB.UnsLabel, string, *relationDB.UnsLabel](labelsPos, func(e *relationDB.UnsLabel) (ok bool, k string, v *relationDB.UnsLabel) {
		return true, e.LabelName, e
	})
	insertList := make([]*relationDB.UnsLabel, 0, len(labels))
	updateTime := time.Now()
	for _, label := range labels {
		po, has := labelMap[label]
		if !has {
			po = &relationDB.UnsLabel{ID: common.NextId(), LabelName: label}
			insertList = append(insertList, po)
			po.CreateAt = updateTime
		}
	}
	if len(insertList) > 0 {
		err = l.labelMapper.MultiInsert(db, insertList)
	}
	return err
}
func (l *UnsLabelService) MakeLabel(ctx context.Context, unsLabels []bo.UnsLabels, createTime time.Time) (rs []*relationDB.UnsLabel, er error) {
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
	db := relationDB.GetDb(ctx)
	if len(resetUnsIds) > 0 {
		er := l.labelRefMapper.DeleteByUnsIds(db, resetUnsIds)
		if er != nil {
			return nil, er
		}
	}
	labels := base.MapKeys[string, []bo.UnsLabels](labelUnsMap)
	allLabels := make([]*relationDB.UnsLabel, 0, len(labels))
	var saveLabels []*relationDB.UnsLabel
	saveLabelRef := make([]*relationDB.UnsLabelRef, 0, len(labels))
	if len(labels) > 0 {
		existLabels, er := l.labelMapper.FindByNames(db, labels)
		if er != nil {
			return nil, er
		}
		existLabelMap := base.MapArrayToMap[*relationDB.UnsLabel, string, *relationDB.UnsLabel](existLabels, func(e *relationDB.UnsLabel) (ok bool, k string, v *relationDB.UnsLabel) {
			return true, e.LabelName, e
		})
		for _, label := range labels {
			existLabel := existLabelMap[label]
			labelPo := existLabel
			if labelPo == nil {
				labelPo = &relationDB.UnsLabel{LabelName: label}
			}
			if existLabel == nil {
				// 新增标签
				labelPo.ID = common.NextId()
				if saveLabels == nil {
					saveLabels = make([]*relationDB.UnsLabel, 0, len(labels))
				}
				saveLabels = append(saveLabels, labelPo)
				allLabels = append(allLabels, labelPo)
			}
			labelId := labelPo.ID
			unsLabelsList := labelUnsMap[label]
			for _, ul := range unsLabelsList {
				ul.SetLabelId(label, labelId)
			}
			saveLabelRef = append(saveLabelRef, base.Map[bo.UnsLabels, *relationDB.UnsLabelRef](unsLabelsList, func(e bo.UnsLabels) *relationDB.UnsLabelRef {
				return &relationDB.UnsLabelRef{LabelID: labelId, UnsID: e.UnsId()}
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
