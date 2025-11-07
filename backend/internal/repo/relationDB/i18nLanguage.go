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

type I18nLanguageRepo struct {
    db *gorm.DB
}

func NewI18nLanguageRepo(in any) *I18nLanguageRepo {
    return &I18nLanguageRepo{db: stores.GetCommonConn(in)}
}

type I18nLanguageFilter struct {
    OnlyForEnabled *bool // 仅返回启用语言
}

func (p I18nLanguageRepo) fmtFilter(ctx context.Context, f I18nLanguageFilter) *gorm.DB {
    db := p.db.WithContext(ctx)
    if f.OnlyForEnabled != nil && *f.OnlyForEnabled {
        db = db.Where("has_used = ?", true)
    }
    return db
}

func (p I18nLanguageRepo) Insert(ctx context.Context, data *SuposI18nLanguage) error {
    result := p.db.WithContext(ctx).Create(data)
    return stores.ErrFmt(result.Error)
}

func (p I18nLanguageRepo) FindOneByFilter(ctx context.Context, f I18nLanguageFilter) (*SuposI18nLanguage, error) {
    var result SuposI18nLanguage
    db := p.fmtFilter(ctx, f)
    err := db.First(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return &result, nil
}
func (p I18nLanguageRepo) FindByFilter(ctx context.Context, f I18nLanguageFilter, page *stores.PageInfo) ([]*SuposI18nLanguage, error) {
    var results []*SuposI18nLanguage
    db := p.fmtFilter(ctx, f).Model(&SuposI18nLanguage{})
    if page != nil {
        db = page.ToGorm(db)
    }
    err := db.Find(&results).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return results, nil
}

func (p I18nLanguageRepo) CountByFilter(ctx context.Context, f I18nLanguageFilter) (size int64, err error) {
    db := p.fmtFilter(ctx, f).Model(&SuposI18nLanguage{})
    err = db.Count(&size).Error
    return size, stores.ErrFmt(err)
}

func (p I18nLanguageRepo) Update(ctx context.Context, data *SuposI18nLanguage) error {
    err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
    return stores.ErrFmt(err)
}

func (p I18nLanguageRepo) DeleteByFilter(ctx context.Context, f I18nLanguageFilter) error {
    db := p.fmtFilter(ctx, f)
    err := db.Delete(&SuposI18nLanguage{}).Error
    return stores.ErrFmt(err)
}

func (p I18nLanguageRepo) Delete(ctx context.Context, id int64) error {
    err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SuposI18nLanguage{}).Error
    return stores.ErrFmt(err)
}
func (p I18nLanguageRepo) FindOne(ctx context.Context, id int64) (*SuposI18nLanguage, error) {
    var result SuposI18nLanguage
    err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p I18nLanguageRepo) MultiInsert(ctx context.Context, data []*SuposI18nLanguage) error {
    err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SuposI18nLanguage{}).Create(data).Error
    return stores.ErrFmt(err)
}

func (d I18nLanguageRepo) UpdateWithField(ctx context.Context, f I18nLanguageFilter, updates map[string]any) error {
    db := d.fmtFilter(ctx, f)
    err := db.Model(&SuposI18nLanguage{}).Updates(updates).Error
    return stores.ErrFmt(err)
}
