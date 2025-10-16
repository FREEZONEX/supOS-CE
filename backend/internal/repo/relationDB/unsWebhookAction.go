package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsWebhookActionRepo struct {
	db *gorm.DB
}

func NewUnsWebhookActionRepo(in any) *UnsWebhookActionRepo {
	return &UnsWebhookActionRepo{db: stores.GetCommonConn(in)}
}

type UnsWebhookActionFilter struct {
	//todo 添加过滤字段
}

func (p UnsWebhookActionRepo) fmtFilter(ctx context.Context, f UnsWebhookActionFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsWebhookActionRepo) Insert(ctx context.Context, data *UnsWebhookAction) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsWebhookActionRepo) FindOneByFilter(ctx context.Context, f UnsWebhookActionFilter) (*UnsWebhookAction, error) {
	var result UnsWebhookAction
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsWebhookActionRepo) FindByFilter(ctx context.Context, f UnsWebhookActionFilter, page *stores.PageInfo) ([]*UnsWebhookAction, error) {
	var results []*UnsWebhookAction
	db := p.fmtFilter(ctx, f).Model(&UnsWebhookAction{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsWebhookActionRepo) CountByFilter(ctx context.Context, f UnsWebhookActionFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsWebhookAction{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsWebhookActionRepo) Update(ctx context.Context, data *UnsWebhookAction) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsWebhookActionRepo) DeleteByFilter(ctx context.Context, f UnsWebhookActionFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsWebhookAction{}).Error
	return stores.ErrFmt(err)
}

func (p UnsWebhookActionRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsWebhookAction{}).Error
	return stores.ErrFmt(err)
}
func (p UnsWebhookActionRepo) FindOne(ctx context.Context, id int64) (*UnsWebhookAction, error) {
	var result UnsWebhookAction
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsWebhookActionRepo) MultiInsert(ctx context.Context, data []*UnsWebhookAction) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsWebhookAction{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsWebhookActionRepo) UpdateWithField(ctx context.Context, f UnsWebhookActionFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsWebhookAction{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
