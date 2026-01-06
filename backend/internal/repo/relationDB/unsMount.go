package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsMountRepo struct {
	db *gorm.DB
}

func NewUnsMountRepo(ctx context.Context) *UnsMountRepo {
	return &UnsMountRepo{db: GetDb(ctx)}
}

type UnsMountFilter struct {
	//todo 添加过滤字段
}

func (p UnsMountRepo) fmtFilter(ctx context.Context, f UnsMountFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsMountRepo) Insert(ctx context.Context, data *UnsMount) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsMountRepo) FindOneByFilter(ctx context.Context, f UnsMountFilter) (*UnsMount, error) {
	var result UnsMount
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsMountRepo) FindByFilter(ctx context.Context, f UnsMountFilter, page *stores.PageInfo) ([]*UnsMount, error) {
	var results []*UnsMount
	db := p.fmtFilter(ctx, f).Model(&UnsMount{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsMountRepo) CountByFilter(ctx context.Context, f UnsMountFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsMount{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsMountRepo) Update(ctx context.Context, data *UnsMount) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsMountRepo) DeleteByFilter(ctx context.Context, f UnsMountFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsMount{}).Error
	return stores.ErrFmt(err)
}

func (p UnsMountRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsMount{}).Error
	return stores.ErrFmt(err)
}

// FindAll returns all mount records
func (p UnsMountRepo) FindAll(ctx context.Context) ([]*UnsMount, error) {
	var results []*UnsMount
	err := p.db.WithContext(ctx).Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsMountRepo) FindOne(ctx context.Context, id int64) (*UnsMount, error) {
	var result UnsMount
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsMountRepo) MultiInsert(ctx context.Context, data []*UnsMount) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsMount{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsMountRepo) UpdateWithField(ctx context.Context, f UnsMountFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsMount{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
