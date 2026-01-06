package relationDB

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsPersonConfigRepo struct {
	db *gorm.DB
}

func NewUnsPersonConfigRepo(in any) *UnsPersonConfigRepo {
	return &UnsPersonConfigRepo{db: stores.GetCommonConn(in)}
}

type UnsPersonConfigFilter struct {
	IDs     []int64
	UserID  string
	UserIDs []string
}

func (p UnsPersonConfigRepo) fmtFilter(ctx context.Context, f UnsPersonConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.IDs) > 0 {
		db = db.Where("id IN ?", f.IDs)
	}
	if uid := strings.TrimSpace(f.UserID); uid != "" {
		db = db.Where("user_id = ?", uid)
	}
	if len(f.UserIDs) > 0 {
		db = db.Where("user_id IN ?", f.UserIDs)
	}
	return db
}

func (p UnsPersonConfigRepo) Insert(ctx context.Context, data *UnsPersonConfig) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsPersonConfigRepo) FindOneByFilter(ctx context.Context, f UnsPersonConfigFilter) (*UnsPersonConfig, error) {
	var result UnsPersonConfig
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsPersonConfigRepo) FindByFilter(ctx context.Context, f UnsPersonConfigFilter, page *stores.PageInfo) ([]*UnsPersonConfig, error) {
	var results []*UnsPersonConfig
	db := p.fmtFilter(ctx, f).Model(&UnsPersonConfig{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsPersonConfigRepo) CountByFilter(ctx context.Context, f UnsPersonConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsPersonConfig{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsPersonConfigRepo) Update(ctx context.Context, data *UnsPersonConfig) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsPersonConfigRepo) DeleteByFilter(ctx context.Context, f UnsPersonConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsPersonConfig{}).Error
	return stores.ErrFmt(err)
}

func (p UnsPersonConfigRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsPersonConfig{}).Error
	return stores.ErrFmt(err)
}
func (p UnsPersonConfigRepo) FindOne(ctx context.Context, id int64) (*UnsPersonConfig, error) {
	var result UnsPersonConfig
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsPersonConfigRepo) MultiInsert(ctx context.Context, data []*UnsPersonConfig) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsPersonConfig{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsPersonConfigRepo) UpdateWithField(ctx context.Context, f UnsPersonConfigFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsPersonConfig{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
