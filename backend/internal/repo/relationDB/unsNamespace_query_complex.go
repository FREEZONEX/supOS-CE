package relationDB

import (
	"backend/internal/common/dto"
	"backend/share/base"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

func (p UnsNamespaceRepo) ListByConditions(db *gorm.DB, f *dto.UnsSearchCondition) (results []*UnsNamespace, err error) {
	db = p.model(db)
	sql := &base.StringBuilder{}
	sql.Grow(512)
	sql.Append("SELECT a.* FROM uns_namespace a ")
	if f.LabelName != "" {
		sql.Append("LEFT JOIN uns_label_ref b ON a.id = b.uns_id   LEFT JOIN uns_label c ON c.id = b.label_id ")
	}
	filterByConditions(sql, f)
	err = db.Raw(sql.String()).Find(&results).Error
	return
}
func (p UnsNamespaceRepo) PageListByConditions(db *gorm.DB, f *dto.UnsSearchCondition, page *stores.PageInfo) (pageRs dto.PageResultDTO[*UnsNamespace], err error) {
	db = p.model(db)
	db = page.ToGorm(db)
	sql := &base.StringBuilder{}
	sql.Grow(512)
	sql.Append(`
		SELECT a.* , t.name AS template_name,t.alias AS template_alias,
			(SELECT COALESCE(string_agg(l.label_name,','),'') FROM uns_label_ref r LEFT JOIN uns_label l ON r.label_id = l.id WHERE r.uns_id = a.id) AS labels
			FROM uns_namespace a
			LEFT JOIN uns_namespace t ON a.model_id = t.id
			LEFT JOIN uns_label_ref b ON a.id = b.uns_id
			LEFT JOIN uns_label c ON c.id = b.label_id
    `)
	filterByConditions(sql, f)
	var sqlStr = sql.String()
	err = db.Raw(sqlStr).Find(&pageRs.Data).Error
	if err != nil {
		return pageRs, stores.ErrFmt(err)
	}
	countErr := db.Raw("select count(*) from (" + sqlStr + ")").Count(&pageRs.Total).Error
	if countErr != nil {
		return pageRs, stores.ErrFmt(countErr)
	}
	pageRs.PageNo = page.Page
	pageRs.PageSize = page.Size
	return
}
func (p UnsNamespaceRepo) PageListTemplates(db *gorm.DB, f dto.TemplateQueryVo, pageNo, pageSize int64, searchCount *int64) (pageRs []*UnsNamespace, err error) {
	db = p.model(db)
	db = db.Select([]string{"id", "name", "path", "alias", "description"}).Where("path_type=1 AND data_type=0 AND status=1 ")
	if f.Key != "" {
		db = db.Where(fmt.Sprintf(" name LIKE '%%%s%%'  ", escapeLikePattern(escapeSQL(f.Key))))
	}
	if subscribeEnable := f.SubscribeEnable; subscribeEnable != nil {
		if *subscribeEnable {
			db = db.Where(" with_flags&256 >0 ")
		} else {
			db = db.Where(" with_flags&256 =0 ")
		}
	}
	if searchCount != nil {
		countErr := db.Count(searchCount).Error
		if countErr != nil {
			return pageRs, stores.ErrFmt(countErr)
		}
	}
	page := &stores.PageInfo{Page: pageNo, Size: pageSize, Orders: []stores.OrderBy{{Field: "id"}}}
	db = page.ToGorm(db)
	err = db.Find(&pageRs).Error
	if err != nil {
		return pageRs, stores.ErrFmt(err)
	}
	return
}
func (p UnsNamespaceRepo) PageListByTemplateId(db *gorm.DB, templateId int64, pageNo, pageSize int64, searchCount *int64) (pageRs []*UnsNamespace, err error) {
	db = p.model(db)
	db = db.Select([]string{"id", "path_type", "name", "path"}).Where("model_id=? AND status=1", templateId)
	if searchCount != nil {
		countErr := db.Count(searchCount).Error
		if countErr != nil {
			err = stores.ErrFmt(countErr)
			return
		} else if *searchCount == 0 {
			return nil, nil
		}
	}
	page := &stores.PageInfo{Page: pageNo, Size: pageSize, Orders: []stores.OrderBy{{Field: "id"}}}
	db = page.ToGorm(db)
	err = db.Find(&pageRs).Error
	if err != nil {
		return pageRs, stores.ErrFmt(err)
	}
	return
}
func filterByConditions(s *base.StringBuilder, f *dto.UnsSearchCondition) {
	s.Append("WHERE a.status = 1 AND (a.data_type != 5 OR a.data_type IS NULL)")

	if f.ParentId != nil {
		if *f.ParentId == 0 {
			s.Append(" AND a.parent_id IS NULL")
		} else {
			s.Append(fmt.Sprintf(" AND a.parent_id = %d", *f.ParentId))
		}
	}

	if f.PathType != nil {
		s.Append("AND a.path_type =").Int(int(*f.PathType))
	} else {
		s.Append(" AND (a.path_type = 0 OR a.path_type = 2)")
	}

	if f.Keyword != "" {
		kw := escapeLikePattern(escapeSQL(f.Keyword))
		s.Append(fmt.Sprintf(
			" AND (a.path ILIKE '%%%s%%' OR a.alias LIKE '%%%s%%')",
			kw, kw))
	}
	if f.Alias != "" {
		s.Append(" AND a.alias = '").Append(escapeSQL(f.Alias)).Append("'")
	}
	if f.ParentAlias == "NULL" {
		s.Append(" AND (a.parent_alias IS NULL OR a.parent_alias = '')")
	} else if f.ParentAlias != "" {
		s.Append(" AND a.parent_alias = '").Append(fmt.Sprintf("'%s'", escapeSQL(f.ParentAlias))).Append("'")
	}
	if len(f.ParentAliasList) > 0 {
		s.Append(" AND a.parent_alias IN (")
		for i, alias := range f.AliasList {
			if i > 0 {
				s.Append(",")
			}
			s.Append("'").Append(escapeSQL(alias)).Append("'")
		}
		s.Append(")")
	}
	if len(f.AliasList) > 0 {
		s.Append(" AND a.alias IN (")
		for i, alias := range f.AliasList {
			if i > 0 {
				s.Append(",")
			}
			s.Append("'").Append(escapeSQL(alias)).Append("'")
		}
		s.Append(")")
	}

	if f.Path != "" {
		s.Append(fmt.Sprintf(" AND a.path LIKE '%%%s%%'", escapeSQL(f.Path)))
	}
	if len(f.PathList) > 0 {
		s.Append(" AND a.path IN (")
		for i, path := range f.PathList {
			if i > 0 {
				s.Append(",")
			}
			s.Append("'").Append(escapeSQL(path)).Append("'")
		}
		s.Append(")")
	}
	if f.LayRec != "" {
		s.Append(fmt.Sprintf(" AND a.lay_rec LIKE '%%%s%%'", escapeSQL(f.LayRec)))
	}
	if f.Name != "" {
		s.Append(" AND a.name = '").Append(escapeSQL(f.Name)).Append("'")
	}
	if f.DisplayName != "" {
		s.Append(fmt.Sprintf(" AND a.display_name LIKE '%%%s%%'", escapeSQL(f.DisplayName)))
	}
	if f.Description != "" {
		s.Append(fmt.Sprintf(" AND a.description LIKE '%%%s%%'", escapeSQL(f.Description)))
	}
	if dataType := f.DataType; dataType != nil {
		s.Append(" AND a.data_type = ").Int(*dataType)
	}
	if templateName := f.TemplateName; templateName != "" {
		s.Append(" AND (a.path_type = 1 AND a.path = '").Append(escapeSQL(templateName)).Append("')")
	}
	if templateId := f.TemplateID; templateId != nil {
		s.Append("a.model_id = ").Long(*templateId)
	}
	if updateStartTime := f.UpdateStartTime; updateStartTime != nil {
		s.Append(fmt.Sprintf(" AND a.update_at >= '%s'::timestamp", *updateStartTime))
	}
	if updateEndTime := f.UpdateEndTime; updateEndTime != nil {
		s.Append(fmt.Sprintf(" AND a.update_at < '%s'::timestamp", *updateEndTime))
	}
	if createStartTime := f.CreateStartTime; createStartTime != nil {
		s.Append(fmt.Sprintf(" AND a.update_at >= '%s'::timestamp", *createStartTime))
	}
	if createEndTime := f.CreateEndTime; createEndTime != nil {
		s.Append(fmt.Sprintf(" AND a.update_at < '%s'::timestamp", *createEndTime))
	}
	if labelName := f.LabelName; labelName != "" {
		s.Append(fmt.Sprintf(" AND c.label_name LIKE '%%%s%%'", escapeSQL(labelName)))
	}
	if len(f.Extend) > 0 {
		for k, v := range f.Extend {
			s.Append(" AND a.extend->> '").Append(escapeSQL(k)).Append("' = '").Append(escapeSQL(fmt.Sprint(v))).Append("'")
		}
	}
	s.Append(" ORDER BY a.id ASC")
}
func escapeSQL(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}
