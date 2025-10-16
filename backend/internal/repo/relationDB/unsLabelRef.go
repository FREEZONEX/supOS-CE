package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsLabelRefRepo struct {
	db *gorm.DB
}

func NewUnsLabelRefRepo(in any) *UnsLabelRefRepo {
	return &UnsLabelRefRepo{db: stores.GetCommonConn(in)}
}

type UnsLabelRefFilter struct {
	//todo 添加过滤字段
}

func (p UnsLabelRefRepo) fmtFilter(ctx context.Context, f UnsLabelRefFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsLabelRefRepo) Insert(ctx context.Context, data *UnsLabelRef) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsLabelRefRepo) FindOneByFilter(ctx context.Context, f UnsLabelRefFilter) (*UnsLabelRef, error) {
	var result UnsLabelRef
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsLabelRefRepo) FindByFilter(ctx context.Context, f UnsLabelRefFilter, page *stores.PageInfo) ([]*UnsLabelRef, error) {
	var results []*UnsLabelRef
	db := p.fmtFilter(ctx, f).Model(&UnsLabelRef{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsLabelRefRepo) CountByFilter(ctx context.Context, f UnsLabelRefFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsLabelRef{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsLabelRefRepo) Update(ctx context.Context, data *UnsLabelRef) error {
	// 组合主键，直接Save
	err := p.db.WithContext(ctx).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsLabelRefRepo) DeleteByFilter(ctx context.Context, f UnsLabelRefFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsLabelRef{}).Error
	return stores.ErrFmt(err)
}

func (p UnsLabelRefRepo) FindOne(ctx context.Context, labelID int64, unsID int64) (*UnsLabelRef, error) {
	var result UnsLabelRef
	err := p.db.WithContext(ctx).Where("label_id = ? AND uns_id = ?", labelID, unsID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsLabelRefRepo) MultiInsert(ctx context.Context, data []*UnsLabelRef) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsLabelRef{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsLabelRefRepo) UpdateWithField(ctx context.Context, f UnsLabelRefFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsLabelRef{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
