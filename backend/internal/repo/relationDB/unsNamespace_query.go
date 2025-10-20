package relationDB

import (
	"backend/internal/common/constants"
	"context"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

func (p UnsNamespaceRepo) dbx(ctx context.Context) *gorm.DB {
	db := p.db.WithContext(ctx).Model(&UnsNamespace{})
	return db
}
func (p UnsNamespaceRepo) ListByAlias(ctx context.Context, alias []string) (results []*UnsNamespace, er error) {
	err := p.dbx(ctx).Where("alias IN ? ", alias).Where("status = ?", 1).Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsNamespaceRepo) GetByAlias(ctx context.Context, alias string) (result *UnsNamespace, err error) {
	var po UnsNamespace
	err = p.dbx(ctx).Where("alias = ? ", alias).Where("status = ?", 1).First(&po).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &po, nil
}
func (p UnsNamespaceRepo) GetAliasByPath(ctx context.Context, path string) (alias string, err error) {
	err = p.dbx(ctx).Select("alias").Where("path = ? ", path).Pluck("alias", &alias).Error
	if err != nil {
		return "", stores.ErrFmt(err)
	}
	return
}

type UnsPathFilter struct {
	TemplateId int64
	Key        string
	PathType   int
	DataTypes  []int
}

func (p UnsNamespaceRepo) CountPaths(ctx context.Context, filter UnsPathFilter) (count int64, er error) {
	db := p.filterPath(ctx, filter)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) filterPath(ctx context.Context, f UnsPathFilter) *gorm.DB {
	db := p.dbx(ctx)
	if f.TemplateId > 0 {
		db = db.Where("model_id = ?", f.TemplateId)
	}
	if f.Key != "" {
		db = db.Where("path iLike ?", f.Key)
	}
	if len(f.DataTypes) > 0 {
		db = db.Where("data_type in ?", f.DataTypes).Where("data_type <> ?", constants.AlarmRuleType)
	}
	db = db.Where("status = ?", 1)
	return db
}

type SimpleUns struct {
	ID       string `gorm:"column:id" json:"id"`
	DataType int    `gorm:"column:data_type" json:"data_type"`
	Alias    string `gorm:"column:alias;not null" json:"alias"`
	Path     string `gorm:"column:path;not null" json:"path"`
}

func (p UnsNamespaceRepo) ListPaths(ctx context.Context, filter UnsPathFilter, page *stores.PageInfo) (results []*SimpleUns, er error) {
	db := p.filterPath(ctx, filter)
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountByDataType(ctx context.Context, key string, dataType int) (count int64, er error) {
	db := p.dbx(ctx)
	db.Where("path_type = ?", 2).Where("data_type = ?", dataType)
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	db = db.Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListFileByIds(ctx context.Context, ids []int64) (results []*UnsNamespace, err error) {
	err = p.dbx(ctx).Where("id in ? ", ids).
		Where("path_type = ?", 2).
		Where("status = ?", 1).
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListFileByTemplateId(ctx context.Context, templateId int64) (results []*UnsNamespace, err error) {
	err = p.dbx(ctx).Where("model_id = ? ", templateId).
		Where("path_type = ?", 2).
		Where("status = ?", 1).
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountNotCalcSeqFiles(ctx context.Context, key string, minNumFields int) (count int64, er error) {
	db := p.dbx(ctx)
	db.Where("path_type = ?", 2).Where("data_type = ?", constants.TimeSequenceType)
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	if minNumFields >= 0 {
		db = db.Where("number_fields", minNumFields)
	}
	db = db.Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListNotCalcSeqFiles(ctx context.Context, key string, minNumFields int, page *stores.PageInfo) (results []*UnsNamespace, err error) {
	db := p.dbx(ctx)
	db = page.ToGorm(db)
	db.Where("path_type = ?", 2).Where("data_type = ?", constants.TimeSequenceType)
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	if minNumFields >= 0 {
		db = db.Where("number_fields", minNumFields)
	}
	db = db.Where("status = ?", 1)
	err = db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountTimeSeriesFiles(ctx context.Context, key string) (count int64, er error) {
	db := p.dbx(ctx)
	db.Where("path_type = ?", 2).Where("data_type in ?", []int64{constants.TimeSequenceType, constants.CalculationRealType})
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	db = db.Where("number_fields > ? ", 0)
	db = db.Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListTimeSeriesFiles(ctx context.Context, key string, page *stores.PageInfo) (results []*UnsNamespace, err error) {
	db := p.dbx(ctx)
	db = page.ToGorm(db)
	db.Where("path_type = ?", 2).Where("data_type in ?", []int64{constants.TimeSequenceType, constants.CalculationRealType})
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	db = db.Where("number_fields > ? ", 0)
	db = db.Where("status = ?", 1)
	err = db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountAlarmRules(ctx context.Context, key string) (count int64, er error) {
	db := p.dbx(ctx)
	db.Where("path_type = ?", 2).Where("data_type = ?", constants.AlarmRuleType)
	if key != "" {
		db = db.Where("(data_path iLike ? OR description like ?)", key, key)
	}
	db = db.Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListAlarmRules(ctx context.Context, key string, page *stores.PageInfo) (results []*UnsNamespace, err error) {
	db := p.dbx(ctx)
	db = page.ToGorm(db)
	db.Where("path_type = ?", 2).Where("data_type = ?", constants.AlarmRuleType)
	if key != "" {
		db = db.Where(db.Where("data_path iLike ?", key).Or("description like ?", key))
	}
	db = db.Where("status = ?", 1)
	err = db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListUnsByIds(ctx context.Context, ids []int64) (results []*UnsPo, err error) {
	query := `
        SELECT a.*, 
               (SELECT COUNT(*) FROM uns_namespace c WHERE c.parent_id = a.id) AS count_direct_children
        FROM uns_namespace a 
        WHERE a.id IN ? AND a.status = 1
    `
	err = p.dbx(ctx).Raw(query, ids).Scan(&results).Error
	return results, err
}
func (p UnsNamespaceRepo) ListInTemplate(ctx context.Context, name string) (results []*UnsNamespace, err error) {
	db := p.dbx(ctx)
	query := db.Where("path_type in ?", []int{0, 2}).
		Where("data_type <> ?", constants.AlarmRuleType).
		Where("model_id IS NOT NULL")
	if name != "" {
		lowerName := strings.ToLower(name)
		query = query.Where(
			"(LOWER(path) LIKE ? OR LOWER(alias) LIKE ?)",
			"%"+lowerName+"%",
			"%"+lowerName+"%",
		)
	}
	err = query.Order("path_type ASC, id ASC").Find(&results).Error
	return results, err
}
func (p UnsNamespaceRepo) CountAllChildrenByLayRec(ctx context.Context, layRec string) (count int64, er error) {
	db := p.dbx(ctx)
	db.Where("path_type = ?", 2).Where("lay_rec like CONCAT('?', '/%')", layRec)
	db = db.Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountDirectChildrenByParentId(ctx context.Context, parentId int64) (count int64, er error) {
	db := p.dbx(ctx).Where("parent_id = ?", parentId).Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}

func (p UnsNamespaceRepo) ListAllEmptyFolder(ctx context.Context) (results []*UnsNamespace, err error) {
	db := p.dbx(ctx)
	query := db.Raw(`select * from ` + TableNameUnsNamespace + `WHERE path_type = 0 and status=1 and (mount_type=0 or mount_type is null) and id NOT IN (
           SELECT DISTINCT parent_id FROM ` + TableNameUnsNamespace + `  WHERE parent_id IS NOT NULL  AND status=1`)
	err = query.Find(&results).Error
	return results, err
}
