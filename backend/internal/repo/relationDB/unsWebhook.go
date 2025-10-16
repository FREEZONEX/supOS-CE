package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsWebhookRepo struct {
	db *gorm.DB
}

func NewUnsWebhookRepo(in any) *UnsWebhookRepo {
	return &UnsWebhookRepo{db: stores.GetCommonConn(in)}
}

type UnsWebhookFilter struct {
	//todo 添加过滤字段
}

func (p UnsWebhookRepo) fmtFilter(ctx context.Context, f UnsWebhookFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsWebhookRepo) Insert(ctx context.Context, data *UnsWebhook) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsWebhookRepo) FindOneByFilter(ctx context.Context, f UnsWebhookFilter) (*UnsWebhook, error) {
	var result UnsWebhook
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsWebhookRepo) FindByFilter(ctx context.Context, f UnsWebhookFilter, page *stores.PageInfo) ([]*UnsWebhook, error) {
	var results []*UnsWebhook
	db := p.fmtFilter(ctx, f).Model(&UnsWebhook{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsWebhookRepo) CountByFilter(ctx context.Context, f UnsWebhookFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsWebhook{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsWebhookRepo) Update(ctx context.Context, data *UnsWebhook) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsWebhookRepo) DeleteByFilter(ctx context.Context, f UnsWebhookFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsWebhook{}).Error
	return stores.ErrFmt(err)
}

func (p UnsWebhookRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsWebhook{}).Error
	return stores.ErrFmt(err)
}
func (p UnsWebhookRepo) FindOne(ctx context.Context, id int64) (*UnsWebhook, error) {
	var result UnsWebhook
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsWebhookRepo) MultiInsert(ctx context.Context, data []*UnsWebhook) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsWebhook{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsWebhookRepo) UpdateWithField(ctx context.Context, f UnsWebhookFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsWebhook{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
