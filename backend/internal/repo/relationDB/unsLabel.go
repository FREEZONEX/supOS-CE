package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsLabelRepo struct {
	db *gorm.DB
}

func NewUnsLabelRepo(in any) *UnsLabelRepo {
	return &UnsLabelRepo{db: stores.GetCommonConn(in)}
}

type UnsLabelFilter struct {
	//todo 添加过滤字段
}

func (p UnsLabelRepo) fmtFilter(ctx context.Context, f UnsLabelFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsLabelRepo) Insert(ctx context.Context, data *UnsLabel) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsLabelRepo) FindOneByFilter(ctx context.Context, f UnsLabelFilter) (*UnsLabel, error) {
	var result UnsLabel
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsLabelRepo) FindByFilter(ctx context.Context, f UnsLabelFilter, page *stores.PageInfo) ([]*UnsLabel, error) {
	var results []*UnsLabel
	db := p.fmtFilter(ctx, f).Model(&UnsLabel{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsLabelRepo) CountByFilter(ctx context.Context, f UnsLabelFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsLabel{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsLabelRepo) Update(ctx context.Context, data *UnsLabel) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsLabelRepo) DeleteByFilter(ctx context.Context, f UnsLabelFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsLabel{}).Error
	return stores.ErrFmt(err)
}

func (p UnsLabelRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsLabel{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRepo) FindOne(ctx context.Context, id int64) (*UnsLabel, error) {
	var result UnsLabel
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsLabelRepo) MultiInsert(ctx context.Context, data []*UnsLabel) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsLabel{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsLabelRepo) UpdateWithField(ctx context.Context, f UnsLabelFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsLabel{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
