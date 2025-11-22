package relationDB

import (
	"backend/internal/common/constants"
	"backend/share/base"
	"context"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p UnsNamespaceRepo) ListByAlias(db *gorm.DB, alias []string) (results []*UnsNamespace, er error) {
	err := p.model(db).Where("alias IN ? ", alias).Where("status = 1").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsNamespaceRepo) ListByIds(db *gorm.DB, ids []int64) (results []*UnsNamespace, er error) {
	err := p.model(db).Where("id IN ? ", ids).Where("status = 1").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

// AllByAlias 忽略逻辑删除标志的按alias查询
func (p UnsNamespaceRepo) AllByAlias(db *gorm.DB, alias []string) (results []*UnsNamespace, er error) {
	err := p.model(db).Where("alias IN ? ", alias).Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

// AllByIds 忽略逻辑删除标志的按Id查询
func (p UnsNamespaceRepo) AllByIds(db *gorm.DB, ids []int64) (results []*UnsNamespace, er error) {
	err := p.model(db).Where("id IN ? ", ids).Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsNamespaceRepo) GetByAlias(db *gorm.DB, alias string) (result *UnsNamespace, err error) {
	var po UnsNamespace
	err = p.model(db).Where("alias = ? ", alias).Where("status = 1").First(&po).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &po, nil
}
func (p UnsNamespaceRepo) GetByPath(db *gorm.DB, path string) (result *UnsNamespace, err error) {
	var po UnsNamespace
	err = p.model(db).
		Where("pathash = hashtext(?)", path).
		Where("path = ? ", path).
		Where("status = 1").First(&po).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &po, nil
}

// ListCategoryFolders 查询分类文件夹
func (p UnsNamespaceRepo) ListCategoryFolders(db *gorm.DB, parentAliasList []string) (results []*UnsNamespace, err error) {
	if len(parentAliasList) == 0 {
		return results, nil
	}
	err = p.model(db).
		Select([]string{"id", "parent_id", "alias", "parent_alias", "data_type"}).
		Where("parent_alias IN ?", parentAliasList).
		Where("data_type > ?", 0).
		Where("path_type = ?", 0).
		Find(&results).Error
	return
}

// ListRootCategoryFolders 查询根分类文件夹
func (p UnsNamespaceRepo) ListRootCategoryFolders(db *gorm.DB) (results []*UnsNamespace, err error) {
	err = p.model(db).
		Select([]string{"id", "parent_id", "alias", "parent_alias", "data_type"}).
		Where("parent_id is null").
		Where("data_type > ?", 0).
		Where("path_type = ?", 0).
		Find(&results).Error
	return
}

type UnsPathFilter struct {
	Key        string
	TemplateId int64
	PathType   int
	DataTypes  []int16
}

type SimpleUns struct {
	ID       string `gorm:"column:id" json:"id"`
	DataType int    `gorm:"column:data_type" json:"data_type"`
	Alias    string `gorm:"column:alias;not null" json:"alias"`
	Path     string `gorm:"column:path;not null" json:"path"`
}

func (p UnsNamespaceRepo) ListPaths(db *gorm.DB, f *UnsPathFilter, page *stores.PageInfo, searchCount *int64) (results []*SimpleUns, er error) {
	db = p.model(db)
	if f.TemplateId > 0 {
		db = db.Where("model_id = ?", f.TemplateId)
	}
	if f.Key != "" {
		db = db.Where("path iLike ?", f.Key)
	}
	if len(f.DataTypes) > 0 {
		db = db.Where("data_type in ?", f.DataTypes).Where("data_type <> ?", constants.AlarmRuleType)
	}
	db = db.Where("id>1000 AND status = 1")
	if searchCount != nil {
		er = db.Count(searchCount).Error
		if er != nil || *searchCount == 0 {
			return
		}
	}
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountByDataType(db *gorm.DB, key string, dataType int) (count int64, er error) {
	db = p.model(db)
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
func (p UnsNamespaceRepo) ListFileByIds(db *gorm.DB, ids []int64) (results []*UnsNamespace, err error) {
	err = p.model(db).Where("id in ? ", ids).
		Where("path_type = ?", 2).
		Where("status = 1").
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListSubTree(db *gorm.DB, layRec string) (results []*UnsNamespace, err error) {
	err = p.model(db).Where("lay_rec like '"+layRec+"/%'").
		Where("status = ?", 1).
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountByParentAliasAndNames(db *gorm.DB, parentAliasAndNames []*UnsNamespace) (results []*UnsNamespace, err error) {
	// 构建VALUES参数
	var sql = &base.StringBuilder{}
	sql.Grow(300 + 80*len(parentAliasAndNames))
	sql.Append(`select  u.parent_id, u."name",max(getIndex(u."path"))+1 as id from (VALUES`)
	for i, data := range parentAliasAndNames {
		if i > 0 {
			sql.Append(",")
		}
		if parentId := data.ParentId; parentId != nil {
			sql.Append(`(`).Long(*parentId)
		} else {
			sql.Append(`(-1`)
		}
		sql.Append(`,'`).Append(escapeSQL(data.Name)).Append(`')`)
	}
	sql.Append(`) AS x(parent_id, name)
	join uns_namespace u on (x.parent_id = u.parent_id or (x.parent_id=-1 and u.parent_id is null)) AND x.name =u.name 
	where u.status =1 group by u.parent_id, u."name"
    `)
	err = p.model(db).Raw(sql.String()).Scan(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	if len(results) > 0 {
		for _, po := range results {
			po.CountExistsSiblings = po.Id
		}
	}
	return
}
func (p UnsNamespaceRepo) ListFileByTemplateId(db *gorm.DB, templateId int64) (results []*UnsNamespace, err error) {
	err = p.model(db).Where("model_id = ? ", templateId).
		Where("path_type = ?", 2).
		Where("status = ?", 1).
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListNotCalcSeqFiles(db *gorm.DB, key string, minNumFields int, page *stores.PageInfo, searchCount *int64) (results []*UnsNamespace, err error) {
	db = p.model(db)
	db.Where("path_type = ?", 2).Where("data_type = ?", constants.TimeSequenceType)
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	if minNumFields >= 0 {
		db = db.Where("number_fields >= ?", minNumFields)
	}
	db = db.Where("status = 1")
	if searchCount != nil {
		err = db.Count(searchCount).Error
		if err != nil || *searchCount == 0 {
			return
		}
	}
	db = page.ToGorm(db)
	err = db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListTimeSeriesFiles(db *gorm.DB, key string, page *stores.PageInfo, searchCount *int64) (results []*UnsNamespace, err error) {
	db = p.model(db)
	db.Where("path_type = ?", 2).Where("data_type in ?", []int16{constants.TimeSequenceType, constants.CalculationRealType})
	if key != "" {
		db = db.Where("path iLike ?", key)
	}
	db = db.Where("number_fields > 0 ")
	db = db.Where("status = 1")
	if searchCount != nil {
		err = db.Count(searchCount).Error
		if err != nil || *searchCount == 0 {
			return
		}
	}
	db = page.ToGorm(db)
	err = db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListAlarmRules(db *gorm.DB, key string, page *stores.PageInfo, searchCount *int64) (results []*UnsNamespace, err error) {
	db = p.model(db)
	db = db.Where("path_type = ?", 2).Where("data_type = ?", constants.AlarmRuleType)
	if key != "" {
		db = db.Where("(data_path like ? OR description like ?)", key, key)
	}
	db = db.Where("status = 1")
	if searchCount != nil {
		err = db.Count(searchCount).Error
		if err != nil || *searchCount == 0 {
			return
		}
	}
	if page != nil {
		db = page.ToGorm(db)
	}
	err = db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) ListByLayRec(db *gorm.DB, layRec string, page *stores.PageInfo, maxId *int64) (results []*UnsNamespace, err error) {
	db = p.model(db).Where("lay_rec like '" + escapeSQL(layRec) + "%'").Where("status=1")
	if maxId != nil {
		err = db.WithContext(context.Background()).Select("MAX(id) as id").Scan(maxId).Error
		if err != nil || *maxId == 0 {
			return
		}
	}
	if page != nil {
		db = page.ToGorm(db)
	}
	err = db.Find(&results).Error
	return
}

func (p UnsNamespaceRepo) ExistsForbiddenFiles(db *gorm.DB, layRec string, pDataType int16) bool {
	db = p.model(db)
	var rs UnsNamespace
	err := db.Select([]string{"id"}).Where("lay_rec like '"+escapeSQL(layRec)+"%'").
		Where("(path_type=0 and data_type <> ?) OR (path_type=2 and parent_data_type <> ?)", pDataType, pDataType).
		Where("status=1").First(&rs).Error
	return err == nil
}
func (p UnsNamespaceRepo) ListByLayRecs(db *gorm.DB, layRecs []string, page *stores.PageInfo) (results []*UnsNamespace, err error) {
	db = p.model(db)
	db = page.ToGorm(db)
	sql := &base.StringBuilder{}
	sql.Grow(80 * len(layRecs))
	sql.Append("select * from ").Append(TableNameUnsNamespace).Append(" WHERE ")
	sql.Append("( ")
	for i, layRec := range layRecs {
		if i > 0 {
			sql.Append(" OR ")
		}
		sql.Append("lay_rec like '").Append(layRec).Append("%'")
	}
	sql.Append(" ) and status=1")
	err = db.Raw(sql.String()).Find(&results).Error
	return
}
func (p UnsNamespaceRepo) ListByTemplateId(db *gorm.DB, templateId int64, page *stores.PageInfo) (results []*UnsNamespace, err error) {
	db = p.model(db)
	if page != nil {
		db = page.ToGorm(db)
	}
	err = db.Where("model_id =?", templateId).Where("status=1").Find(&results).Error
	return
}
func (p UnsNamespaceRepo) ListByTemplateIds(db *gorm.DB, templateIds []int64, page *stores.PageInfo) (results []*UnsNamespace, err error) {
	db = p.model(db)
	if page != nil {
		db = page.ToGorm(db)
	}
	err = db.Where("model_id in ?", templateIds).Where("status=1").Find(&results).Error
	return
}
func (p UnsNamespaceRepo) ListUnsByIds(db *gorm.DB, ids []int64) (results []*UnsPo, err error) {
	query := `
        SELECT a.*, 
               (SELECT COUNT(*) FROM uns_namespace c WHERE c.parent_id = a.id) AS count_direct_children
        FROM uns_namespace a 
        WHERE a.id IN ? AND a.status = 1
    `
	err = p.model(db).Raw(query, ids).Scan(&results).Error
	return results, err
}

func (p UnsNamespaceRepo) ListInTemplate(db *gorm.DB, name string) (results []*UnsNamespace, err error) {
	db = p.model(db)
	query := db.Where("path_type in ?", []int{0, 2}).
		Where("data_type <> ?", constants.AlarmRuleType).
		Where("model_id IS NOT NULL")
	if name != "" {
		lowerName := "%" + strings.ToLower(escapeSQL(name)) + "%"
		query = query.Where(
			"(LOWER(path) LIKE ? OR LOWER(alias) LIKE ?)",
			lowerName,
			lowerName,
		)
	}
	err = query.Order("path_type ASC, id ASC").Find(&results).Error
	return results, err
}
func (p UnsNamespaceRepo) CountAllChildrenByLayRec(db *gorm.DB, layRec string) (count int64, er error) {
	db = p.model(db)
	db.Where("path_type = ?", 2).Where("lay_rec like CONCAT('?', '/%')", layRec)
	db = db.Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) CountDirectChildrenByParentId(db *gorm.DB, parentId int64) (count int64, er error) {
	db = p.model(db).Where("parent_id = ?", parentId).Where("status = ?", 1)
	err := db.Count(&count).Error
	if err != nil {
		return -1, stores.ErrFmt(err)
	}
	return
}

func (p UnsNamespaceRepo) ListAllEmptyFolder(db *gorm.DB) (results []*UnsNamespace, err error) {
	db = p.model(db)
	query := db.Raw(`select * from ` + TableNameUnsNamespace + `WHERE path_type = 0 and status=1 and (mount_type=0 or mount_type is null) and id NOT IN (
           SELECT DISTINCT parent_id FROM ` + TableNameUnsNamespace + `  WHERE parent_id IS NOT NULL  AND status=1`)
	err = query.Find(&results).Error
	return results, err
}

func (p UnsNamespaceRepo) ListLabeledUnsByKeyword(db *gorm.DB, keyword string) (results []*UnsNamespace, err error) {
	query := db.Table("uns_namespace n").
		Where("id in(select distinct uns_id from uns_label_ref ulr )").
		Where("n.path_type = ?", 2).
		Where("n.status =1") // 软删除过滤

	if kw := strings.TrimSpace(keyword); kw != "" {
		likeKeyword := "%" + strings.ToLower(escapeSQL(kw)) + "%"
		query = query.Where("(LOWER(n.path) LIKE ? OR LOWER(n.alias) LIKE ?)", likeKeyword, likeKeyword)
	}

	err = query.Order(clause.OrderByColumn{Column: clause.Column{Name: "n.id"}}).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get uns by keyword: %w", err)
	}

	return results, nil
}
func (p UnsNamespaceRepo) PageListByLabel(db *gorm.DB, labelID int64, pageNo, pageSize int64, searchCount *int64) (unsList []*UnsNamespace, err error) {
	db = p.model(db)
	page := &stores.PageInfo{Page: pageNo, Size: pageSize, Orders: []stores.OrderBy{{Field: "id", Sort: stores.OrderAsc}}}

	// 基础查询
	baseQuery := db.
		Joins("JOIN uns_label_ref rf ON uns_namespace.id = rf.uns_id").
		Where("rf.label_id = ? AND uns_namespace.status=1", labelID)

	// COUNT 查询（不包含 SELECT 子句）
	if searchCount != nil {
		countQuery := baseQuery.Session(&gorm.Session{})
		err = countQuery.Count(searchCount).Error
		if err != nil || *searchCount == 0 {
			return
		}
	}

	// 数据查询（包含 SELECT 子句）
	dataQuery := baseQuery.
		Select("uns_namespace.*").
		Scopes(page.ToGorm)

	err = dataQuery.Find(&unsList).Error
	return unsList, err
}
