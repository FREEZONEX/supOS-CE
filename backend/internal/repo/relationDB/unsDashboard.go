package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsDashboardRepo struct {
	db *gorm.DB
}

func NewUnsDashboardRepo(in any) *UnsDashboardRepo {
	return &UnsDashboardRepo{db: stores.GetCommonConn(in)}
}

type UnsDashboardFilter struct {
	//todo 添加过滤字段
}

func (p UnsDashboardRepo) fmtFilter(ctx context.Context, f UnsDashboardFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsDashboardRepo) Insert(ctx context.Context, data *UnsDashboard) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsDashboardRepo) FindOneByFilter(ctx context.Context, f UnsDashboardFilter) (*UnsDashboard, error) {
	var result UnsDashboard
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsDashboardRepo) FindByFilter(ctx context.Context, f UnsDashboardFilter, page *stores.PageInfo) ([]*UnsDashboard, error) {
	var results []*UnsDashboard
	db := p.fmtFilter(ctx, f).Model(&UnsDashboard{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsDashboardRepo) CountByFilter(ctx context.Context, f UnsDashboardFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsDashboard{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsDashboardRepo) Update(ctx context.Context, data *UnsDashboard) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsDashboardRepo) DeleteByFilter(ctx context.Context, f UnsDashboardFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsDashboard{}).Error
	return stores.ErrFmt(err)
}

func (p UnsDashboardRepo) Delete(ctx context.Context, id string) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsDashboard{}).Error
	return stores.ErrFmt(err)
}
func (p UnsDashboardRepo) FindOne(ctx context.Context, id string) (*UnsDashboard, error) {
	var result UnsDashboard
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsDashboardRepo) MultiInsert(ctx context.Context, data []*UnsDashboard) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsDashboard{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsDashboardRepo) UpdateWithField(ctx context.Context, f UnsDashboardFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsDashboard{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
