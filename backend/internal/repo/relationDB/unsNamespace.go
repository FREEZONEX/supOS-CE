package relationDB

import (
	"context"
	"strings"

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
}

func NewUnsNamespaceRepo() UnsNamespaceRepo {
	return UnsNamespaceRepo{}
}

func GetDb(ctx context.Context) *gorm.DB {
	if connObj := ctx.Value("db"); connObj != nil {
		if db, is := connObj.(*gorm.DB); is {
			return db
		}
	}
	return stores.GetCommonConn(ctx).Debug()
}
func SetDb(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "db", db)
}

type UnsNamespaceFilter struct {
	//todo 添加过滤字段
}

func (p UnsNamespaceRepo) model(db *gorm.DB) *gorm.DB {
	return db.Model(&UnsNamespace{})
}

func (p UnsNamespaceRepo) Insert(db *gorm.DB, data *UnsNamespace) error {
	result := p.model(db).Create(data)
	return stores.ErrFmt(result.Error)
}

// 批量插入记录
func (p UnsNamespaceRepo) MultiInsert(db *gorm.DB, data []*UnsNamespace) error {
	err := p.model(db).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, 1000).Error
	return stores.ErrFmt(err)
}
func (p UnsNamespaceRepo) MultiUpdate(db *gorm.DB, data []*UnsNamespace) (err error) {
	err = p.model(db).Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(data, 1000).Error
	return stores.ErrFmt(err)
	/*	_, isTransaction := db.Statement.ConnPool.(gorm.TxCommitter)
		db = p.model(db)
		if isTransaction {
			for _, item := range data {
				if item.ID == 0 {
					continue
				}
				err = db.Omit("id", "created_at").Updates(item).Error
				if err != nil {
					break
				}
			}
		} else {
			err = db.Transaction(func(tx *gorm.DB) error {
				for _, item := range data {
					if item.ID == 0 {
						continue // 记录日志或返回错误
					}
					if dbErr := tx.Model(&UnsNamespace{}).Omit("id", "created_at").Updates(item).Error; dbErr != nil {
						return fmt.Errorf("failed to update record ID %d: %w", item.ID, dbErr)
					}
				}
				return nil
			})
		}
		return stores.ErrFmt(err)*/
}
func (p UnsNamespaceRepo) FindOneByFilter(db *gorm.DB, f UnsNamespaceFilter) (*UnsNamespace, error) {
	var result UnsNamespace
	err := p.model(db).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsNamespaceRepo) FindByFilter(db *gorm.DB, f UnsNamespaceFilter, page *stores.PageInfo) ([]*UnsNamespace, error) {
	var results []*UnsNamespace
	db = p.model(db)
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsNamespaceRepo) CountByFilter(db *gorm.DB, f UnsNamespaceFilter) (size int64, err error) {
	err = p.model(db).Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) Update(db *gorm.DB, data *UnsNamespace) error {
	err := p.model(db).Where("id = ?", data.Id).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) DeleteByFilter(db *gorm.DB, f UnsNamespaceFilter) error {
	err := p.model(db).Delete(&UnsNamespace{}).Error
	return stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) Delete(db *gorm.DB, id int64) error {
	err := p.model(db).Where("id = ?", id).Delete(&UnsNamespace{}).Error
	return stores.ErrFmt(err)
}
func (p UnsNamespaceRepo) SelectById(db *gorm.DB, id int64) (*UnsNamespace, error) {
	var result UnsNamespace
	err := p.model(db).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	} else if result.Id == 0 {
		return nil, nil
	}
	return &result, nil
}
func (p UnsNamespaceRepo) SelectByIds(db *gorm.DB, ids []int64) (results []*UnsNamespace, err error) {
	err = p.model(db).Where("id IN ?", ids).First(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsNamespaceRepo) FindOneByAlias(db *gorm.DB, alias string) (*UnsNamespace, error) {
	if alias == "" {
		return nil, stores.ErrFmt(gorm.ErrRecordNotFound)
	}
	var result UnsNamespace
	err := p.model(db).Where("alias = ?", alias).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p UnsNamespaceRepo) UpdateWithField(db *gorm.DB, f UnsNamespaceFilter, updates map[string]any) error {
	err := p.model(db).Updates(updates).Error
	return stores.ErrFmt(err)
}

func escapeLikePattern(input string) string {
	if input == "" {
		return input
	}
	input = strings.ReplaceAll(input, `\`, `\\`)
	input = strings.ReplaceAll(input, `%`, `\%`)
	input = strings.ReplaceAll(input, `_`, `\_`)
	return input
}

func (p UnsNamespaceRepo) ListAliasByBase(db *gorm.DB, base string) ([]string, error) {
	if base == "" {
		return nil, nil
	}
	escaped := escapeLikePattern(base)
	pattern := escaped + "-%"
	var aliases []string
	err := p.model(db).
		Model(&UnsNamespace{}).
		Select("alias").
		Where("(alias = ? OR alias LIKE ? )", base, pattern).
		Where("status=1").
		Pluck("alias", &aliases).Error
	return aliases, stores.ErrFmt(err)
}
