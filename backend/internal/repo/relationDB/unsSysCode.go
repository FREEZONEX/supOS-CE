package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsSysCodeRepo struct {
	db *gorm.DB
}

func NewUnsSysCodeRepo(in any) *UnsSysCodeRepo {
	return &UnsSysCodeRepo{db: stores.GetCommonConn(in)}
}

type UnsSysCodeFilter struct {
	//todo 添加过滤字段
}

func (p UnsSysCodeRepo) fmtFilter(ctx context.Context, f UnsSysCodeFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsSysCodeRepo) Insert(ctx context.Context, data *UnsSysCode) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsSysCodeRepo) FindOneByFilter(ctx context.Context, f UnsSysCodeFilter) (*UnsSysCode, error) {
	var result UnsSysCode
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsSysCodeRepo) FindByFilter(ctx context.Context, f UnsSysCodeFilter, page *stores.PageInfo) ([]*UnsSysCode, error) {
	var results []*UnsSysCode
	db := p.fmtFilter(ctx, f).Model(&UnsSysCode{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsSysCodeRepo) CountByFilter(ctx context.Context, f UnsSysCodeFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsSysCode{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsSysCodeRepo) Update(ctx context.Context, data *UnsSysCode) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsSysCodeRepo) DeleteByFilter(ctx context.Context, f UnsSysCodeFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsSysCode{}).Error
	return stores.ErrFmt(err)
}

func (p UnsSysCodeRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsSysCode{}).Error
	return stores.ErrFmt(err)
}
func (p UnsSysCodeRepo) FindOne(ctx context.Context, id int64) (*UnsSysCode, error) {
	var result UnsSysCode
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsSysCodeRepo) MultiInsert(ctx context.Context, data []*UnsSysCode) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsSysCode{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsSysCodeRepo) UpdateWithField(ctx context.Context, f UnsSysCodeFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsSysCode{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
