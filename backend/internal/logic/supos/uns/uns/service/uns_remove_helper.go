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
	"encoding/csv"
	"io"
	"sync"
	"time"

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
		tar, _ := r.unsMapper.SelectById(db, singleId)
		if tar == nil {
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
			list, _ := r.unsMapper.ListByAlias(db, aliasList)
			if len(list) > 0 {
				unsPos = append(unsPos, list...)
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
		err := r.removeFolders(ctx, req, folders)
		if err != nil {
			return nil, err
		}
	}

	if files := paramGroups[constants.PathTypeFile]; len(files) > 0 {
		for _, fs := range base.Partition(files, 1000) {
			err := r.deleteAndSendEvent(ctx, req, fs)
			if err != nil {
				return nil, err
			}
		}
	}

	if templates := paramGroups[constants.PathTypeTemplate]; len(templates) > 0 {
		err := r.removeTemplates(ctx, req, templates)
		if err != nil {
			return nil, err
		}
	}
	return &types.RemoveResult{BaseResult: types.BaseResult{Code: 200, Msg: "ok"}}, nil
}

func (r *UnsRemoveService) removeTemplates(ctx context.Context, req types.RemoveUnsOptions, templates []*dao.UnsNamespace) error {
	var files = make([]*dao.UnsNamespace, 0, 128)
	var folders = make([]*dao.UnsNamespace, 0, 64)
	for _, templateIds := range base.Partition(base.Map(templates, getId), 1000) {
		// 创建管道
		reader, writer := io.Pipe()
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer func() {
				_ = writer.Close()
				wg.Done()
			}()
			delEr := r.unsMapper.ExportCsvByTemplateIds(ctx, templateIds, writer)
			if delEr != nil {
				r.log.Error("Del exportByTemplate failed:", delEr)
			}
		}()
		go func() {
			defer func() {
				_ = reader.Close()
				wg.Done()
			}()
			// 读取 CSV 表头
			csvReader := csv.NewReader(reader)
			headers, err := csvReader.Read()
			if err != nil {
				r.log.Error("del exportHeader ByTemplate failed:", err)
				return
			}
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				unsPO := r.unsMapper.Csv2Model(headers, record)
				switch unsPO.PathType {
				case constants.PathTypeFile:
					files = append(files, unsPO)
					if len(files) >= 1000 {
						delEr := r.deleteAndSendEvent(ctx, req, files)
						if delEr != nil {
							r.log.Error("del exportByTemplate failed:", delEr)
						}
						files = files[:0]
					}
				case constants.PathTypeDir:
					folders = append(folders, unsPO)
				}
			}

		}()
		wg.Wait()
	}
	files = append(files, templates...)
	err := r.deleteAndSendEventWithCall(ctx, req, files, func(db *gorm.DB) (er error) {
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
	return err
}

func (r *UnsRemoveService) removeFolders(
	ctx context.Context,
	req types.RemoveUnsOptions,
	folders []*dao.UnsNamespace,
) error {
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
	for _, lay := range base.Partition(layRecs, 500) {
		// 创建管道
		reader, writer := io.Pipe()
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer func() {
				_ = writer.Close()
				wg.Done()
			}()
			delEr := r.unsMapper.ExportCsvByLayRecAndIds(ctx, lay, nil, writer, false) //layRec降序排列，保证目录是最后一个被删除
			if delEr != nil {
				r.log.Error("Del exportByFolder failed:", delEr)
			}
		}()
		go func() {
			defer func() {
				_ = reader.Close()
				defer wg.Done()
			}()
			// 读取 CSV 表头
			csvReader := csv.NewReader(reader)
			headers, err := csvReader.Read()
			if err != nil {
				r.log.Error("del exportHeader ByFolder failed:", err)
				return
			}
			batch := make([]*dao.UnsNamespace, 0, 1000)
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				unsPO := r.unsMapper.Csv2Model(headers, record)
				batch = append(batch, unsPO)
				if len(batch) >= 1000 {
					delEr := r.deleteAndSendEvent(ctx, req, batch)
					batch = batch[:0]
					if delEr != nil {
						r.log.Error("delByFolder failed:", delEr)
					}
				}
			}
			if len(batch) > 0 {
				delEr := r.deleteAndSendEvent(ctx, req, batch)
				if delEr != nil {
					r.log.Error("DelByFolder failed:", delEr)
				}
			}
		}()
		wg.Wait()
	}
	return nil
}
func getId(e *dao.UnsNamespace) int64 {
	return e.Id
}
func getPathType(e *dao.UnsNamespace) int16 {
	return e.PathType
}
func (r *UnsRemoveService) deleteAndSendEvent(ctx context.Context, req types.RemoveUnsOptions, list []*dao.UnsNamespace) (err error) {
	return r.deleteAndSendEventWithCall(ctx, req, list, nil)
}
func (r *UnsRemoveService) deleteAndSendEventWithCall(ctx context.Context, req types.RemoveUnsOptions, list []*dao.UnsNamespace, callback func(db *gorm.DB) error) (er error) {
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
		delEvent := event.NewRemoveTopicsEvent(ctx, time.Now(), withFlow, withDashboard,
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
