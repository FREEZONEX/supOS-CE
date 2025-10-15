package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsAlarmsHandlerRepo struct {
	db *gorm.DB
}

func NewUnsAlarmsHandlerRepo(in any) *UnsAlarmsHandlerRepo {
	return &UnsAlarmsHandlerRepo{db: stores.GetCommonConn(in)}
}

type UnsAlarmsHandlerFilter struct {
	//todo 添加过滤字段
}

func (p UnsAlarmsHandlerRepo) fmtFilter(ctx context.Context, f UnsAlarmsHandlerFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsAlarmsHandlerRepo) Insert(ctx context.Context, data *UnsAlarmsHandler) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsAlarmsHandlerRepo) FindOneByFilter(ctx context.Context, f UnsAlarmsHandlerFilter) (*UnsAlarmsHandler, error) {
	var result UnsAlarmsHandler
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsAlarmsHandlerRepo) FindByFilter(ctx context.Context, f UnsAlarmsHandlerFilter, page *stores.PageInfo) ([]*UnsAlarmsHandler, error) {
	var results []*UnsAlarmsHandler
	db := p.fmtFilter(ctx, f).Model(&UnsAlarmsHandler{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsAlarmsHandlerRepo) CountByFilter(ctx context.Context, f UnsAlarmsHandlerFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsAlarmsHandler{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsAlarmsHandlerRepo) Update(ctx context.Context, data *UnsAlarmsHandler) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsAlarmsHandlerRepo) DeleteByFilter(ctx context.Context, f UnsAlarmsHandlerFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsAlarmsHandler{}).Error
	return stores.ErrFmt(err)
}

func (p UnsAlarmsHandlerRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsAlarmsHandler{}).Error
	return stores.ErrFmt(err)
}
func (p UnsAlarmsHandlerRepo) FindOne(ctx context.Context, id int64) (*UnsAlarmsHandler, error) {
	var result UnsAlarmsHandler
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsAlarmsHandlerRepo) MultiInsert(ctx context.Context, data []*UnsAlarmsHandler) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsAlarmsHandler{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsAlarmsHandlerRepo) UpdateWithField(ctx context.Context, f UnsAlarmsHandlerFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsAlarmsHandler{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
