package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsDashboardTopRecodeRepo struct {
	db *gorm.DB
}

func NewUnsDashboardTopRecodeRepo(in any) *UnsDashboardTopRecodeRepo {
	return &UnsDashboardTopRecodeRepo{db: stores.GetCommonConn(in)}
}

type UnsDashboardTopRecodeFilter struct {
	//todo 添加过滤字段
}

func (p UnsDashboardTopRecodeRepo) fmtFilter(ctx context.Context, f UnsDashboardTopRecodeFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsDashboardTopRecodeRepo) Insert(ctx context.Context, data *UnsDashboardTopRecode) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsDashboardTopRecodeRepo) FindOneByFilter(ctx context.Context, f UnsDashboardTopRecodeFilter) (*UnsDashboardTopRecode, error) {
	var result UnsDashboardTopRecode
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsDashboardTopRecodeRepo) FindByFilter(ctx context.Context, f UnsDashboardTopRecodeFilter, page *stores.PageInfo) ([]*UnsDashboardTopRecode, error) {
	var results []*UnsDashboardTopRecode
	db := p.fmtFilter(ctx, f).Model(&UnsDashboardTopRecode{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsDashboardTopRecodeRepo) CountByFilter(ctx context.Context, f UnsDashboardTopRecodeFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsDashboardTopRecode{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsDashboardTopRecodeRepo) Update(ctx context.Context, data *UnsDashboardTopRecode) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsDashboardTopRecodeRepo) DeleteByFilter(ctx context.Context, f UnsDashboardTopRecodeFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsDashboardTopRecode{}).Error
	return stores.ErrFmt(err)
}

func (p UnsDashboardTopRecodeRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsDashboardTopRecode{}).Error
	return stores.ErrFmt(err)
}
func (p UnsDashboardTopRecodeRepo) FindOne(ctx context.Context, id int64) (*UnsDashboardTopRecode, error) {
	var result UnsDashboardTopRecode
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsDashboardTopRecodeRepo) MultiInsert(ctx context.Context, data []*UnsDashboardTopRecode) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsDashboardTopRecode{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsDashboardTopRecodeRepo) UpdateWithField(ctx context.Context, f UnsDashboardTopRecodeFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsDashboardTopRecode{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
