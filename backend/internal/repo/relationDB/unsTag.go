package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsTagRepo struct {
	db *gorm.DB
}

func NewUnsTagRepo(in any) *UnsTagRepo {
	return &UnsTagRepo{db: stores.GetCommonConn(in)}
}

type UnsTagFilter struct {
	//todo 添加过滤字段
}

func (p UnsTagRepo) fmtFilter(ctx context.Context, f UnsTagFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsTagRepo) Insert(ctx context.Context, data *UnsTag) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsTagRepo) FindOneByFilter(ctx context.Context, f UnsTagFilter) (*UnsTag, error) {
	var result UnsTag
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsTagRepo) FindByFilter(ctx context.Context, f UnsTagFilter, page *stores.PageInfo) ([]*UnsTag, error) {
	var results []*UnsTag
	db := p.fmtFilter(ctx, f).Model(&UnsTag{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsTagRepo) CountByFilter(ctx context.Context, f UnsTagFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsTag{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsTagRepo) Update(ctx context.Context, data *UnsTag) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsTagRepo) DeleteByFilter(ctx context.Context, f UnsTagFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsTag{}).Error
	return stores.ErrFmt(err)
}

func (p UnsTagRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsTag{}).Error
	return stores.ErrFmt(err)
}
func (p UnsTagRepo) FindOne(ctx context.Context, id int64) (*UnsTag, error) {
	var result UnsTag
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsTagRepo) MultiInsert(ctx context.Context, data []*UnsTag) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsTag{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsTagRepo) UpdateWithField(ctx context.Context, f UnsTagFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsTag{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
