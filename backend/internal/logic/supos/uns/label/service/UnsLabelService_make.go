package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"time"
)

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
func (l *UnsLabelService) ClearAndMakeLabels(ctx context.Context, unsId int64, labelList []*types.LabelVo) error {
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.SelectById(db, unsId)
	if uns == nil {
		return nil
	} else if er != nil {
		return er
	}
	ctx = dao.SetDb(ctx, db)
	labelIdNameMap := make(map[int64]string)
	uns.LabelIds = labelIdNameMap
	er = l.labelMapper.DeleteRefByUnsId(db, unsId) //先清空UNS下所有标签
	if er != nil {
		return er
	}

	if len(labelList) > 0 {
		var noNames = base.MapArrayToMap[*types.LabelVo, int64, *types.LabelVo](labelList, func(e *types.LabelVo) (ok bool, k int64, v *types.LabelVo) {
			if e.ID != 0 && e.LabelName == "" {
				ok, k, v = true, e.ID, e
			}
			return
		})

		if len(noNames) > 0 {
			ids := base.MapKeys[int64](noNames)
			labels, er := l.labelMapper.ListByIds(db, ids)
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
		addLabels := make([]string, 0, len(labelList))
		for _, labelVo := range labelList {
			lid := labelVo.ID
			var ref *dao.UnsLabelRef
			if lid != 0 {
				ref = &dao.UnsLabelRef{LabelID: lid, UnsID: unsId}
				labelIdNameMap[ref.LabelID] = labelVo.LabelName
			} else if labelName := labelVo.LabelName; len(labelName) > 0 {
				addLabels = append(addLabels, labelName)
			}
		}
		if len(addLabels) > 0 {
			idMap, err := l.CreateBatch(ctx, addLabels)
			if err != nil {
				return err
			} else if len(idMap) > 0 {
				for name, id := range idMap {
					labelIdNameMap[id] = name
				}
			}
		}
		if len(labelIdNameMap) > 0 {
			er = l.labelRefMapper.SaveOrIgnore(db, base.Map[int64, *dao.UnsLabelRef](base.MapKeys(labelIdNameMap), func(labelId int64) *dao.UnsLabelRef {
				return &dao.UnsLabelRef{LabelID: labelId, UnsID: unsId}
			}))
			if er != nil {
				return er
			}
		}
	}
	uns.UpdateAt = time.Now()
	er = l.unsMapper.Update(db, uns)
	return er
}

func (l *UnsLabelService) MakeSingleLabel(ctx context.Context, req *types.MakeSingleLabelReq) (resp *types.BaseResult, err error) {
	db := dao.GetDb(ctx)
	var labelPo *dao.UnsLabel
	unsId, labelId := req.UnsId, req.LabelId
	labelPo, err = l.labelMapper.SelectById(db, labelId)
	resp = &types.BaseResult{}
	if err != nil || labelPo == nil {
		resp.Code, resp.Msg = 400, I18nUtils.GetMessage("uns.label.not.exists")
		return
	}
	ref, _ := l.labelRefMapper.FindOne(db, labelId, unsId)
	if ref == nil {
		err = l.unsMapper.LinkLabelOnUns(db, unsId, labelId, labelPo.LabelName, time.Now())
		if err != nil {
			return
		}
		uns, _ := l.unsMapper.SelectById(db, unsId)
		if uns == nil {
			resp.Code, resp.Msg = 400, I18nUtils.GetMessage("uns.file.not.exist")
			return
		}
		err = l.labelRefMapper.SaveOrIgnore(db, []*dao.UnsLabelRef{{LabelID: labelId, UnsID: unsId}})
	}
	resp.Code, resp.Msg = 200, "OK"
	return
}
