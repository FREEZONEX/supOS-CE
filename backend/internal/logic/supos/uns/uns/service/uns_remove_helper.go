package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"time"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

func defaultFalse(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func (r *UnsRemoveService) removeModelOrInstance(ctx context.Context, singleId int64, req types.BatchRemoveUnsDto) (resp *types.RemoveResult, err error) {
	db := dao.GetDb(ctx)
	resp = &types.RemoveResult{BaseResult: types.BaseResult{Code: 200, Msg: "ok"}}
	var unsPos []*dao.UnsNamespace
	withFlow, withDashboard := defaultFalse(req.WithFlow), defaultFalse(req.WithDashboard)
	if singleId > 0 {
		tar, err := r.unsMapper.SelectById(db, singleId)
		if err != nil {
			resp.Code = 500
			return resp, err
		} else if tar == nil {
			resp.Code = 400
			resp.Msg = I18nUtils.GetMessage("uns.folder.or.file.not.found")
			return resp, err
		}
		unsPos = []*dao.UnsNamespace{tar}
	} else if len(req.AliasList) == 0 {
		resp.Msg = "NoUnsInParams"
		return
	} else {
		unsPos = make([]*dao.UnsNamespace, 0, len(req.AliasList))
		for _, aliasList := range base.Partition(req.AliasList, 1000) {
			list, er := r.unsMapper.ListByAlias(db, aliasList)
			if len(list) > 0 {
				unsPos = append(unsPos, list...)
			} else if er != nil {
				err = er
				resp.Code = 500
				return resp, err
			}
		}
		if len(unsPos) == 0 {
			resp.Code = 400
			resp.Msg = I18nUtils.GetMessage("uns.folder.or.file.not.found")
			return resp, err
		}
	}
	paramGroups := base.GroupBy[*dao.UnsNamespace, int16](unsPos, func(e *dao.UnsNamespace) int16 {
		return e.PathType
	})
	folders := paramGroups[constants.PathTypeDir]
	if len(folders) > 0 {
		onlyRemoveChild := false
		if mrc := req.OnlyRemoveChild; mrc != nil {
			onlyRemoveChild = *mrc
		}
		var layRecs []string
		if onlyRemoveChild {
			layRecs = base.Map[*dao.UnsNamespace, string](folders, func(e *dao.UnsNamespace) string {
				return e.LayRec + "/"
			})
		} else {
			layRecs = base.Map[*dao.UnsNamespace, string](folders, func(e *dao.UnsNamespace) string {
				return e.LayRec
			})
		}
		loop := true
		for _, lay := range base.Partition(layRecs, 500) {
			page := &stores.PageInfo{Page: 1, Size: 1000, Orders: []stores.OrderBy{{Field: "id", Sort: stores.OrderAsc}}}

			for loop {
				rs, er := r.unsMapper.ListByLayRecs(db, lay, page)
				if er != nil {
					return nil, er
				} else if len(rs) == 0 {
					break
				}
				page.Page++
				unsGroups := base.MapAndGroupBy[*dao.UnsNamespace, *types.CreateTopicDto, int16](rs, func(e *dao.UnsNamespace) (int16, *types.CreateTopicDto) {
					return e.PathType, UnsConverter.Po2Dto(e)
				})
				err = db.Transaction(func(tx *gorm.DB) error {
					er = r.unsMapper.LogicDeleteByIds(tx, base.Map[*dao.UnsNamespace, int64](rs, func(e *dao.UnsNamespace) int64 {
						return e.Id
					}))
					if er == nil {
						delEvent := event.NewRemoveTopicsEvent(dao.SetDb(ctx, tx), time.Now(), withFlow, withDashboard,
							unsGroups[constants.PathTypeFile],
							unsGroups[constants.PathTypeTemplate],
							unsGroups[constants.PathTypeDir],
						)
						er = spring.PublishEvent(delEvent)
						if er != nil {
							resp.Msg = "RemoveEventErr"
						}
					} else {
						resp.Msg = "DbErr"
					}
					return er
				})
				if err != nil {
					resp.Code = 500
					loop = false
					break
				}
			}
			if !loop {
				break
			}
		}
	}
	files := paramGroups[constants.PathTypeFile]
	if len(files) > 0 && err == nil {
		for _, fs := range base.Partition(files, 1000) {
			err = db.Transaction(func(tx *gorm.DB) error {
				er := r.unsMapper.LogicDeleteByIds(tx, base.Map[*dao.UnsNamespace, int64](fs, func(e *dao.UnsNamespace) int64 {
					return e.Id
				}))
				if er == nil {
					delEvent := event.NewRemoveTopicsEvent(dao.SetDb(ctx, tx), time.Now(), withFlow, withDashboard,
						UnsConverter.Po2Dtos(fs), nil, nil,
					)
					er = spring.PublishEvent(delEvent)
					if er != nil {
						resp.Msg = "RemoveEventErr"
					}
				} else {
					resp.Msg = "DbErr"
				}
				return er
			})
			if err != nil {
				resp.Code = 500
				break
			}
		}
	}
	return
}
