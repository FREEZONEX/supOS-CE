package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsAlarmsDatumRepo struct {
	db *gorm.DB
}

func NewUnsAlarmsDatumRepo(in context.Context) *UnsAlarmsDatumRepo {
	return &UnsAlarmsDatumRepo{db: GetDb(in)}
}

type UnsAlarmsDatumFilter struct {
	//todo 添加过滤字段
}

func (p UnsAlarmsDatumRepo) fmtFilter(ctx context.Context, f UnsAlarmsDatumFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsAlarmsDatumRepo) Insert(ctx context.Context, data *UnsAlarmsDatum) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsAlarmsDatumRepo) FindOneByFilter(ctx context.Context, f UnsAlarmsDatumFilter) (*UnsAlarmsDatum, error) {
	var result UnsAlarmsDatum
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// Count returns the total count of alarm records
func (p UnsAlarmsDatumRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&UnsAlarmsDatum{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}

func (p UnsAlarmsDatumRepo) FindByFilter(ctx context.Context, f UnsAlarmsDatumFilter, page *stores.PageInfo) ([]*UnsAlarmsDatum, error) {
	var results []*UnsAlarmsDatum
	db := p.fmtFilter(ctx, f).Model(&UnsAlarmsDatum{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsAlarmsDatumRepo) CountByFilter(ctx context.Context, f UnsAlarmsDatumFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsAlarmsDatum{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsAlarmsDatumRepo) Update(ctx context.Context, data *UnsAlarmsDatum) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsAlarmsDatumRepo) DeleteByFilter(ctx context.Context, f UnsAlarmsDatumFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsAlarmsDatum{}).Error
	return stores.ErrFmt(err)
}

func (p UnsAlarmsDatumRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsAlarmsDatum{}).Error
	return stores.ErrFmt(err)
}
func (p UnsAlarmsDatumRepo) FindOne(ctx context.Context, id int64) (*UnsAlarmsDatum, error) {
	var result UnsAlarmsDatum
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsAlarmsDatumRepo) MultiInsert(ctx context.Context, data []*UnsAlarmsDatum) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsAlarmsDatum{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsAlarmsDatumRepo) UpdateWithField(ctx context.Context, f UnsAlarmsDatumFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsAlarmsDatum{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
