package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/event"
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

func (l *UnsLabelService) Create(ctx context.Context, req *types.CreateLabelReq) (resp *types.CreateLabelResp, err error) {
	// 参数校验
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.Parameter.WithMsg("labelName不能为空")
	}

	resp = &types.CreateLabelResp{}
	db := dao.GetDb(ctx)
	countExits, er := l.labelMapper.CountByName(db, req.Name)
	if er != nil {
		resp.Code = 500
		resp.Msg = er.Error()
		err = er
		return
	} else if countExits > 0 {
		resp.Code = 400
		resp.Msg = I18nUtils.GetMessage("uns.label.already.exists")
		return
	}
	// 写入数据库
	data := &dao.UnsLabel{
		ID:        common.NextId(),
		LabelName: req.Name,
	}
	if err = l.labelMapper.Insert(db, data); err != nil {
		return nil, err
	}
	resp.Code = 200
	resp.Data = &types.LabelVo{
		ID: data.ID,
	}
	return
}
func (l *UnsLabelService) Delete(ctx context.Context, req *types.WithID) (resp *types.BaseResult, err error) {
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

	return &types.BaseResult{Code: 200}, err
}

func (l *UnsLabelService) Update(ctx context.Context, req *types.UpdateLabelReq) (resp *types.BaseResult, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	//if strings.TrimSpace(req.LabelName) == "" {
	//	return nil, errors.Parameter.WithMsg("labelName不能为空")
	//}
	if req.SubscribeFrequency == "" && req.SubscribeEnable == nil && req.LabelName == "" {
		return nil, errors.Parameter.WithMsg("没有什么可以更新")
	}
	db := dao.GetDb(ctx)
	item, _ := l.labelMapper.FindOne(db, req.ID)
	if item == nil {
		return &types.BaseResult{Code: 400, Msg: I18nUtils.GetMessage("uns.label.not.exists")}, nil
	}
	items, er := l.labelMapper.FindByNames(db, []string{req.LabelName})
	if er != nil {
		return nil, er
	} else if len(items) > 1 || (len(items) == 1 && items[0].ID != item.ID) {
		return &types.BaseResult{Code: 400, Msg: I18nUtils.GetMessage("uns.label.already.exists")}, nil
	}
	item.UpdateAt = time.Now()
	item.LabelName = req.LabelName
	item.SubscribeFrequency = req.SubscribeFrequency
	if enable := req.SubscribeEnable; enable != nil {
		flags := base.P2v(item.WithFlags)
		if *enable {
			flags |= constants.UnsFlagWithSubscribeEnable
		} else {
			flags &= ^constants.UnsFlagWithSubscribeEnable
		}
		item.WithFlags = &flags
	}

	if err = l.labelMapper.Update(db, item); err != nil {
		return nil, err
	}
	if len(req.LabelName) > 0 {
		labelId, labelName := req.ID, req.LabelName
		err = l.unsMapper.UpdateUnsLabelNames(db, labelId, labelName) // 更新uns冗余的标签 id->name 键值对
	}
	resp = &types.BaseResult{Code: 200}
	return
}
func (l *UnsLabelService) UpdateSubscribe(ctx context.Context, req *types.UpdateLabelSubscribeReq) (resp *types.BaseResult, err error) {
	return l.Update(ctx, &types.UpdateLabelReq{
		ID:                 req.ID,
		SubscribeEnable:    req.SubscribeEnable,
		SubscribeFrequency: req.SubscribeFrequency,
	})
}

func (l *UnsLabelService) CreateBatch(ctx context.Context, labels []string) (labelIdMap map[string]int64, err error) {
	if len(labels) == 0 {
		return
	}
	db := dao.GetDb(ctx)
	labelsPos, _ := l.labelMapper.FindByNames(db, labels)
	labelMap := base.MapArrayToMap[*dao.UnsLabel, string, *dao.UnsLabel](labelsPos, func(e *dao.UnsLabel) (ok bool, k string, v *dao.UnsLabel) {
		return true, e.LabelName, e
	})
	insertList := make([]*dao.UnsLabel, 0, len(labels))
	updateTime := time.Now()
	labelIdMap = make(map[string]int64, len(labels))
	for _, label := range labels {
		po, has := labelMap[label]
		if !has {
			po = &dao.UnsLabel{ID: common.NextId(), LabelName: label, CreateAt: updateTime}
			insertList = append(insertList, po)
		}
		labelIdMap[po.LabelName] = po.ID
	}
	if len(insertList) > 0 {
		err = l.labelMapper.MultiInsert(db, insertList)
	}
	return labelIdMap, err
}
func (l *UnsLabelService) CancelLabel(ctx context.Context, unsId int64, labelIds []int64) error {
	if len(labelIds) == 0 {
		return nil
	}
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.SelectById(db, unsId)
	if uns == nil {
		return nil
	} else if er != nil {
		return er
	}
	er = l.labelRefMapper.DeleteByUnsIdAndLabelIds(db, uns.Id, labelIds)
	if er != nil {
		return er
	}
	leftLabels, er := l.labelMapper.ListByUnsId(db, unsId)
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
	er = l.unsMapper.Update(db, updatePo)
	if er != nil {
		return er
	}
	return nil
}
func (l *UnsLabelService) CancelLabelByNames(ctx context.Context, unsAlias string, labelNames []string) error {
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.GetByAlias(db, unsAlias)
	if uns == nil {
		return errors.NewCodeError(400, I18nUtils.GetMessage("uns.file.not.exist"))
	} else if er != nil {
		return er
	}
	labelList, er := l.labelMapper.FindByNames(db, labelNames)
	if len(labelList) == 0 {
		return nil
	} else if er != nil {
		return nil
	}
	labelIds := base.Map[*dao.UnsLabel, int64](labelList, func(e *dao.UnsLabel) int64 {
		return e.ID
	})
	ctx = dao.SetDb(ctx, db)
	return l.CancelLabel(ctx, uns.Id, labelIds)
}

// OnEventRemoveTopicsEvent 处理UNS 删除事件
func (l *UnsLabelService) OnEventRemoveTopicsEvent(event *event.RemoveTopicsEvent) (er error) {
	labelIds := base.FilterAndFlatMap(event.Topics, func(e *types.CreateTopicDto) (vs []int64, ok bool) {
		if labelIDs := e.GetLabelIds(); len(labelIDs) > 0 {
			vs, ok = base.MapKeys(labelIDs), true
		}
		return
	})
	if len(labelIds) == 0 {
		return nil
	}
	db := dao.GetDb(context.Background())
	for _, partLabelIds := range base.Partition(labelIds, 1000) {
		er = l.labelRefMapper.DeleteByLabelIds(db, partLabelIds)
		if er != nil {
			return er
		}
	}
	return
}
