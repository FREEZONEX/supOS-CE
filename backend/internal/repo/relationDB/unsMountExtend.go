package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsMountExtendRepo struct {
	db *gorm.DB
}

func NewUnsMountExtendRepo(in any) *UnsMountExtendRepo {
	return &UnsMountExtendRepo{db: stores.GetCommonConn(in)}
}

type UnsMountExtendFilter struct {
	//todo 添加过滤字段
}

func (p UnsMountExtendRepo) fmtFilter(ctx context.Context, f UnsMountExtendFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsMountExtendRepo) Insert(ctx context.Context, data *UnsMountExtend) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsMountExtendRepo) FindOneByFilter(ctx context.Context, f UnsMountExtendFilter) (*UnsMountExtend, error) {
	var result UnsMountExtend
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsMountExtendRepo) FindByFilter(ctx context.Context, f UnsMountExtendFilter, page *stores.PageInfo) ([]*UnsMountExtend, error) {
	var results []*UnsMountExtend
	db := p.fmtFilter(ctx, f).Model(&UnsMountExtend{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsMountExtendRepo) CountByFilter(ctx context.Context, f UnsMountExtendFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsMountExtend{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsMountExtendRepo) Update(ctx context.Context, data *UnsMountExtend) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsMountExtendRepo) DeleteByFilter(ctx context.Context, f UnsMountExtendFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsMountExtend{}).Error
	return stores.ErrFmt(err)
}

func (p UnsMountExtendRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsMountExtend{}).Error
	return stores.ErrFmt(err)
}
func (p UnsMountExtendRepo) FindOne(ctx context.Context, id int64) (*UnsMountExtend, error) {
	var result UnsMountExtend
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsMountExtendRepo) MultiInsert(ctx context.Context, data []*UnsMountExtend) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsMountExtend{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsMountExtendRepo) UpdateWithField(ctx context.Context, f UnsMountExtendFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsMountExtend{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
