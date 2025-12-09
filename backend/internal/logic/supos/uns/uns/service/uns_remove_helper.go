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
func (r *UnsRemoveService) Remove(ctx context.Context, req types.RemoveUnsOptions, unsList []*dao.UnsNamespace) (*types.RemoveResult, error) {
	ctx = dao.SetDb(ctx, dao.GetDb(ctx))
	paramGroups := base.GroupBy(unsList, getPathType)
	if folders := paramGroups[constants.PathTypeDir]; len(folders) > 0 {
		rs, err := r.removeFolders(ctx, req, folders)
		if rs != nil || err != nil {
			return rs, err
		}
	}

	if files := paramGroups[constants.PathTypeFile]; len(files) > 0 {
		for _, fs := range base.Partition(files, 1000) {
			rs, err := r.deleteAndSendEvent(ctx, req, fs)
			if rs != nil || err != nil {
				return rs, err
			}
		}
	}

	if templates := paramGroups[constants.PathTypeTemplate]; len(templates) > 0 {
		rs, err := r.removeTemplates(ctx, req, templates)
		if rs != nil || err != nil {
			return rs, err
		}
	}
	return &types.RemoveResult{BaseResult: types.BaseResult{Code: 200, Msg: "ok"}}, nil
}

func (r *UnsRemoveService) removeTemplates(ctx context.Context, req types.RemoveUnsOptions, templates []*dao.UnsNamespace) (*types.RemoveResult, error) {
	db := dao.GetDb(ctx)
	var files = make([]*dao.UnsNamespace, 0, 128)
	var folders = make([]*dao.UnsNamespace, 0, 64)
	for _, templateIds := range base.Partition(base.Map(templates, getId), 1000) {
		page := &stores.PageInfo{Page: 1, Size: 1000, Orders: []stores.OrderBy{{Field: "id"}}}
		for {
			list, er := r.unsMapper.ListByTemplateIds(db, templateIds, page)
			if er != nil || len(list) == 0 {
				break
			}
			page.Page++

			unsGroups := base.GroupBy(list, getPathType)
			if fs := unsGroups[constants.PathTypeFile]; len(fs) > 0 {
				files = append(files, fs...)
			}
			if dirs := unsGroups[constants.PathTypeDir]; len(dirs) > 0 {
				folders = append(folders, dirs...)
			}
			if len(files) >= 1000 {
				rs, err := r.deleteAndSendEvent(ctx, req, files)
				if rs != nil || err != nil {
					return rs, err
				}
				files = files[:]
			}
			if int64(len(list)) < page.Size {
				break
			}
		}
	}
	files = append(files, templates...)
	rs, err := r.deleteAndSendEventWithCall(ctx, req, files, func(db *gorm.DB) (er error) {
		if len(folders) > 0 {
			dirIds := base.Map(folders, getId)
			if len(dirIds) < 1000 {
				_, er = r.unsMapper.UpdateNullTemplateIdByIds(db, dirIds)
			} else {
				for _, partIds := range base.Partition(dirIds, 1000) {
					_, er = r.unsMapper.UpdateNullTemplateIdByIds(db, partIds)
					if er != nil {
						break
					}
				}
			}
		}
		return
	})
	return rs, err
}

func (r *UnsRemoveService) removeFolders(
	ctx context.Context,
	req types.RemoveUnsOptions,
	folders []*dao.UnsNamespace,
) (*types.RemoveResult, error) {
	onlyRemoveChild := base.P2v(req.OnlyRemoveChild)
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
	db := dao.GetDb(ctx)
	for _, lay := range base.Partition(layRecs, 500) {
		page := &stores.PageInfo{Page: 1, Size: 1000, Orders: []stores.OrderBy{{Field: "id", Sort: stores.OrderAsc}}}
		for {
			list, er := r.unsMapper.ListByLayRecs(db, lay, page)
			if er != nil {
				return nil, er
			} else if len(list) == 0 {
				break
			}
			page.Page++
			delRs, delEr := r.deleteAndSendEvent(ctx, req, list)
			if delEr != nil {
				return delRs, delEr
			} else if delRs != nil && delRs.Data != nil {
				return delRs, nil
			}
			if int64(len(list)) < page.Size {
				break
			}
		}
	}
	return nil, nil
}
func getId(e *dao.UnsNamespace) int64 {
	return e.Id
}
func getPathType(e *dao.UnsNamespace) int16 {
	return e.PathType
}
func (r *UnsRemoveService) deleteAndSendEvent(ctx context.Context, req types.RemoveUnsOptions, list []*dao.UnsNamespace) (resp *types.RemoveResult, err error) {
	return r.deleteAndSendEventWithCall(ctx, req, list, nil)
}
func (r *UnsRemoveService) deleteAndSendEventWithCall(ctx context.Context, req types.RemoveUnsOptions, list []*dao.UnsNamespace, callback func(db *gorm.DB) error) (resp *types.RemoveResult, er error) {
	unsGroups := base.MapAndGroupBy[*dao.UnsNamespace, *types.CreateTopicDto, int16](list, func(e *dao.UnsNamespace) (int16, *types.CreateTopicDto) {
		return e.PathType, UnsConverter.Po2Dto(e)
	})
	files := unsGroups[constants.PathTypeFile]
	if len(files) > 0 {
		//TODO 引用检查
	}
	db := dao.GetDb(ctx)
	var tx = db
	withTx := false
	if !dao.IsInTransaction(db) {
		tx = db.Begin()
		withTx = true
	}
	ids := base.Map[*dao.UnsNamespace, int64](list, func(e *dao.UnsNamespace) int64 {
		return e.Id
	})
	if len(ids) <= 1000 {
		er = r.unsMapper.LogicDeleteByIds(tx, ids)
	} else {
		for _, partIds := range base.Partition(ids, 1000) {
			er = r.unsMapper.LogicDeleteByIds(tx, partIds)
		}
	}
	if er == nil && callback != nil {
		er = callback(tx)
	}

	if er == nil {
		withFlow, withDashboard := defaultFalse(req.WithFlow), defaultFalse(req.WithDashboard)
		delEvent := event.NewRemoveTopicsEvent(context.Background(), time.Now(), withFlow, withDashboard,
			files,
			unsGroups[constants.PathTypeTemplate],
			unsGroups[constants.PathTypeDir],
		)
		er = spring.PublishEvent(delEvent)
	}
	if er == nil {
		if withTx {
			tx.Commit()
		}
	} else if withTx {
		r.log.Error("UNS删除回滚:", er)
		tx.Rollback()
	} else {
		r.log.Error("UNS删除失败:", er)
	}
	return
}
