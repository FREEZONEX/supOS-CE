package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsDashboardRefRepo struct {
	db *gorm.DB
}

func NewUnsDashboardRefRepo(in any) *UnsDashboardRefRepo {
	return &UnsDashboardRefRepo{db: stores.GetCommonConn(in)}
}

type UnsDashboardRefFilter struct {
	//todo 添加过滤字段
}

func (p UnsDashboardRefRepo) fmtFilter(ctx context.Context, f UnsDashboardRefFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsDashboardRefRepo) Insert(ctx context.Context, data *UnsDashboardRef) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsDashboardRefRepo) FindOneByFilter(ctx context.Context, f UnsDashboardRefFilter) (*UnsDashboardRef, error) {
	var result UnsDashboardRef
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsDashboardRefRepo) FindByFilter(ctx context.Context, f UnsDashboardRefFilter, page *stores.PageInfo) ([]*UnsDashboardRef, error) {
	var results []*UnsDashboardRef
	db := p.fmtFilter(ctx, f).Model(&UnsDashboardRef{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsDashboardRefRepo) CountByFilter(ctx context.Context, f UnsDashboardRefFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsDashboardRef{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsDashboardRefRepo) Update(ctx context.Context, data *UnsDashboardRef) error {
	// 该表没有单一自增ID，直接Save
	err := p.db.WithContext(ctx).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsDashboardRefRepo) DeleteByFilter(ctx context.Context, f UnsDashboardRefFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsDashboardRef{}).Error
	return stores.ErrFmt(err)
}

func (p UnsDashboardRefRepo) FindOne(ctx context.Context, dashboardID string, unsAlias string) (*UnsDashboardRef, error) {
	var result UnsDashboardRef
	err := p.db.WithContext(ctx).Where("dashboard_id = ? AND uns_alias = ?", dashboardID, unsAlias).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsDashboardRefRepo) MultiInsert(ctx context.Context, data []*UnsDashboardRef) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsDashboardRef{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsDashboardRefRepo) UpdateWithField(ctx context.Context, f UnsDashboardRefFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsDashboardRef{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
