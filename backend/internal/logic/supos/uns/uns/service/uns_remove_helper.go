package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/logic/supos/uns/uns/bo"
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
	ctx = dao.SetDb(ctx, db)
	return r.Remove(ctx, req.RemoveUnsOptions, unsPos)
}
func (r *UnsRemoveService) Remove(ctx context.Context, req types.RemoveUnsOptions, unsList []*dao.UnsNamespace) (resp *types.RemoveResult, err error) {
	withFlow, withDashboard := defaultFalse(req.WithFlow), defaultFalse(req.WithDashboard)
	resp = &types.RemoveResult{BaseResult: types.BaseResult{Code: 200, Msg: "ok"}}

	db := dao.GetDb(ctx)

	paramGroups := base.GroupBy(unsList, func(e *dao.UnsNamespace) int16 {
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
			layRecs = base.Map(folders, func(e *dao.UnsNamespace) string {
				return e.LayRec + "/"
			})
		} else {
			layRecs = base.Map(folders, func(e *dao.UnsNamespace) string {
				return e.LayRec
			})
		}
		for _, lay := range base.Partition(layRecs, 500) {
			page := &stores.PageInfo{Page: 1, Size: 1000, Orders: []stores.OrderBy{{Field: "id", Sort: stores.OrderAsc}}}

			for {
				rs, er := r.unsMapper.ListByLayRecs(db, lay, page)
				if er != nil {
					return nil, er
				} else if len(rs) == 0 {
					break
				}
				page.Page++
				errMsg := r.deleteAndSendEvent(ctx, db, rs, withFlow, withDashboard)
				if len(errMsg) > 0 {
					resp.Code = 500
					resp.Msg = errMsg
					return
				}
			}
		}
	}
	files := paramGroups[constants.PathTypeFile]
	if len(files) > 0 {
		for _, fs := range base.Partition(files, 1000) {
			errMsg := r.deleteAndSendEvent(ctx, db, fs, withFlow, withDashboard)
			if len(errMsg) > 0 {
				resp.Code = 500
				resp.Msg = errMsg
				return
			}
		}
	}
	return
}

func (r *UnsRemoveService) deleteAndSendEvent(ctx context.Context, db *gorm.DB, rs []*dao.UnsNamespace, withFlow bool, withDashboard bool) (erMsg string) {
	unsGroups := base.MapAndGroupBy[*dao.UnsNamespace, bo.UnsInfo, int16](rs, func(e *dao.UnsNamespace) (int16, bo.UnsInfo) {
		return e.PathType, e
	})

	tx := db.Begin()
	er := r.unsMapper.LogicDeleteByIds(tx, base.Map[*dao.UnsNamespace, int64](rs, func(e *dao.UnsNamespace) int64 {
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
			erMsg = "RemoveEventErr:" + er.Error()
		}
	} else {
		erMsg = "DbErr"
	}
	if er == nil {
		tx.Commit()
	} else {
		if len(erMsg) == 0 {
			erMsg = er.Error()
		}
		tx.Rollback()
	}
	return erMsg
}
