package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsSysModuleRepo struct {
	db *gorm.DB
}

func NewUnsSysModuleRepo(in any) *UnsSysModuleRepo {
	return &UnsSysModuleRepo{db: stores.GetCommonConn(in)}
}

type UnsSysModuleFilter struct {
	//todo 添加过滤字段
}

func (p UnsSysModuleRepo) fmtFilter(ctx context.Context, f UnsSysModuleFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsSysModuleRepo) Insert(ctx context.Context, data *UnsSysModule) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsSysModuleRepo) FindOneByFilter(ctx context.Context, f UnsSysModuleFilter) (*UnsSysModule, error) {
	var result UnsSysModule
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsSysModuleRepo) FindByFilter(ctx context.Context, f UnsSysModuleFilter, page *stores.PageInfo) ([]*UnsSysModule, error) {
	var results []*UnsSysModule
	db := p.fmtFilter(ctx, f).Model(&UnsSysModule{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsSysModuleRepo) CountByFilter(ctx context.Context, f UnsSysModuleFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsSysModule{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsSysModuleRepo) Update(ctx context.Context, data *UnsSysModule) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsSysModuleRepo) DeleteByFilter(ctx context.Context, f UnsSysModuleFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsSysModule{}).Error
	return stores.ErrFmt(err)
}

func (p UnsSysModuleRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsSysModule{}).Error
	return stores.ErrFmt(err)
}
func (p UnsSysModuleRepo) FindOne(ctx context.Context, id int64) (*UnsSysModule, error) {
	var result UnsSysModule
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsSysModuleRepo) MultiInsert(ctx context.Context, data []*UnsSysModule) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsSysModule{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsSysModuleRepo) UpdateWithField(ctx context.Context, f UnsSysModuleFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsSysModule{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
