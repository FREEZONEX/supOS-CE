package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsExportRecordRepo struct {
	db *gorm.DB
}

func NewUnsExportRecordRepo(in any) *UnsExportRecordRepo {
	return &UnsExportRecordRepo{db: stores.GetCommonConn(in)}
}

type UnsExportRecordFilter struct {
	//todo 添加过滤字段
}

func (p UnsExportRecordRepo) fmtFilter(ctx context.Context, f UnsExportRecordFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsExportRecordRepo) Insert(ctx context.Context, data *UnsExportRecord) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsExportRecordRepo) FindOneByFilter(ctx context.Context, f UnsExportRecordFilter) (*UnsExportRecord, error) {
	var result UnsExportRecord
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsExportRecordRepo) FindByFilter(ctx context.Context, f UnsExportRecordFilter, page *stores.PageInfo) ([]*UnsExportRecord, error) {
	var results []*UnsExportRecord
	db := p.fmtFilter(ctx, f).Model(&UnsExportRecord{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsExportRecordRepo) CountByFilter(ctx context.Context, f UnsExportRecordFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsExportRecord{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsExportRecordRepo) Update(ctx context.Context, data *UnsExportRecord) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsExportRecordRepo) DeleteByFilter(ctx context.Context, f UnsExportRecordFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsExportRecord{}).Error
	return stores.ErrFmt(err)
}

func (p UnsExportRecordRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsExportRecord{}).Error
	return stores.ErrFmt(err)
}
func (p UnsExportRecordRepo) FindOne(ctx context.Context, id int64) (*UnsExportRecord, error) {
	var result UnsExportRecord
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsExportRecordRepo) MultiInsert(ctx context.Context, data []*UnsExportRecord) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsExportRecord{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsExportRecordRepo) UpdateWithField(ctx context.Context, f UnsExportRecordFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsExportRecord{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
