package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type UnsNamespaceRepo struct {
	db *gorm.DB
}

func NewUnsNamespaceRepo(in any) *UnsNamespaceRepo {
	if in == nil {
		in = stores.GetCommonConn(context.Background()).Debug()
	}
	return &UnsNamespaceRepo{db: stores.GetCommonConn(in)}
}

type UnsNamespaceFilter struct {
	//todo 添加过滤字段
}

func (p UnsNamespaceRepo) fmtFilter(ctx context.Context, f UnsNamespaceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsNamespaceRepo) Insert(ctx context.Context, data *UnsNamespace) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsNamespaceRepo) FindOneByFilter(ctx context.Context, f UnsNamespaceFilter) (*UnsNamespace, error) {
	var result UnsNamespace
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsNamespaceRepo) FindByFilter(ctx context.Context, f UnsNamespaceFilter, page *stores.PageInfo) ([]*UnsNamespace, error) {
	var results []*UnsNamespace
	db := p.fmtFilter(ctx, f).Model(&UnsNamespace{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsNamespaceRepo) CountByFilter(ctx context.Context, f UnsNamespaceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsNamespace{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) Update(ctx context.Context, data *UnsNamespace) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) DeleteByFilter(ctx context.Context, f UnsNamespaceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsNamespace{}).Error
	return stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsNamespace{}).Error
	return stores.ErrFmt(err)
}
func (p UnsNamespaceRepo) FindOne(ctx context.Context, id int64) (*UnsNamespace, error) {
	var result UnsNamespace
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsNamespaceRepo) MultiInsert(ctx context.Context, data []*UnsNamespace) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsNamespace{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsNamespaceRepo) UpdateWithField(ctx context.Context, f UnsNamespaceFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsNamespace{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
