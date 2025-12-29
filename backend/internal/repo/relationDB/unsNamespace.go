package relationDB

import (
	"backend/internal/common/utils/loggerlevel"
	"context"
	"strings"
	"sync"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
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
	db := stores.GetCommonConn(ctx)
	if loggerlevel.IsDebug() {
		db = db.Debug() //打日志
	} else {
		db = db.Session(&gorm.Session{
			Logger: db.Logger.LogMode(logger.Silent), //不打印日志
		})
	}
	return db
}
func SetDb(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "db", db)
}

func IsInTransaction(db *gorm.DB) bool {
	// 比较当前连接池是否与原始连接池相同
	_, isTransaction := db.Statement.ConnPool.(gorm.TxCommitter)
	return isTransaction
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

var updateColumnsInit sync.Once
var updateColumns []string

func (p UnsNamespaceRepo) MultiUpdate(db *gorm.DB, data []*UnsNamespace) (err error) {
	if len(updateColumns) == 0 {
		// 获取模型的所有字段（除了主键）
		updateColumnsInit.Do(func() {
			var model UnsNamespace
			stmt := db.Model(&model).Statement
			stmt.Parse(&model)
			for _, field := range stmt.Schema.Fields {
				// 跳过主键字段
				if field.PrimaryKey {
					continue
				}
				if len(field.DBName) > 0 {
					updateColumns = append(updateColumns, field.DBName)
				}
			}
		})
	}
	err = p.model(db).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}, clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(data, 1000).Error

	return stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) Update(db *gorm.DB, data *UnsNamespace) error {
	// 使用 Updates 方法（自动忽略 nil 字段）
	err := p.model(db).Where("id = ?", data.Id).Where("status=1").Omit("id").Updates(data).Error
	return stores.ErrFmt(err)
}

func (p UnsNamespaceRepo) Delete(db *gorm.DB, id int64) error {
	err := p.model(db).Where("id = ?", id).Delete(&UnsNamespace{}).Error
	return stores.ErrFmt(err)
}
func (p UnsNamespaceRepo) SelectById(db *gorm.DB, id int64) (*UnsNamespace, error) {
	var result UnsNamespace
	err := p.model(db).Where("id = ?", id).Where("status=1").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	} else if result.Id == 0 {
		return nil, nil
	}
	return &result, nil
}
func (p UnsNamespaceRepo) SelectByIds(db *gorm.DB, ids []int64) (results []*UnsNamespace, err error) {
	err = p.model(db).Where("id IN ?", ids).Where("status=1").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: false}).
		Find(&results).Error
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
	err := p.model(db).Where("alias = ?", alias).Where("status=1").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func EscapeLike(s string) string {
	return escapeLikePattern(escapeSQL(s))
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

// PathTypeCount represents the count of UNS by path type
type PathTypeCount struct {
	PathType int   `gorm:"column:path_type"`
	Count    int64 `gorm:"column:count"`
}

// ProtocolCount represents the count of UNS by protocol type
type ProtocolCount struct {
	Protocol string `gorm:"column:protocol_type"`
	Count    int64  `gorm:"column:count"`
}

// CountByPathType returns count grouped by path_type
func (p UnsNamespaceRepo) CountByPathType(db *gorm.DB) ([]PathTypeCount, error) {
	var results []PathTypeCount
	err := p.model(db).
		Model(&UnsNamespace{}).
		Select("path_type, COUNT(*) as count").
		Where("status = 1 and id>1000").
		Group("path_type").
		Find(&results).Error
	return results, stores.ErrFmt(err)
}

// CountByProtocolType returns count grouped by protocol_type
func (p UnsNamespaceRepo) CountByProtocolType(db *gorm.DB) ([]ProtocolCount, error) {
	var results []ProtocolCount
	err := p.model(db).
		Model(&UnsNamespace{}).
		Select("protocol_type, COUNT(*) as count").
		Where("status = 1 AND protocol_type IS NOT NULL AND protocol_type != ''").
		Group("protocol_type").
		Find(&results).Error
	return results, stores.ErrFmt(err)
}
