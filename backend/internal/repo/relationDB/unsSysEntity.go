package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsSysEntityRepo struct {
	db *gorm.DB
}

func NewUnsSysEntityRepo(in any) *UnsSysEntityRepo {
	return &UnsSysEntityRepo{db: stores.GetCommonConn(in)}
}

type UnsSysEntityFilter struct {
	//todo 添加过滤字段
}

func (p UnsSysEntityRepo) fmtFilter(ctx context.Context, f UnsSysEntityFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsSysEntityRepo) Insert(ctx context.Context, data *UnsSysEntity) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsSysEntityRepo) FindOneByFilter(ctx context.Context, f UnsSysEntityFilter) (*UnsSysEntity, error) {
	var result UnsSysEntity
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsSysEntityRepo) FindByFilter(ctx context.Context, f UnsSysEntityFilter, page *stores.PageInfo) ([]*UnsSysEntity, error) {
	var results []*UnsSysEntity
	db := p.fmtFilter(ctx, f).Model(&UnsSysEntity{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsSysEntityRepo) CountByFilter(ctx context.Context, f UnsSysEntityFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsSysEntity{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsSysEntityRepo) Update(ctx context.Context, data *UnsSysEntity) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsSysEntityRepo) DeleteByFilter(ctx context.Context, f UnsSysEntityFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsSysEntity{}).Error
	return stores.ErrFmt(err)
}

func (p UnsSysEntityRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsSysEntity{}).Error
	return stores.ErrFmt(err)
}
func (p UnsSysEntityRepo) FindOne(ctx context.Context, id int64) (*UnsSysEntity, error) {
	var result UnsSysEntity
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsSysEntityRepo) MultiInsert(ctx context.Context, data []*UnsSysEntity) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsSysEntity{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsSysEntityRepo) UpdateWithField(ctx context.Context, f UnsSysEntityFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsSysEntity{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
