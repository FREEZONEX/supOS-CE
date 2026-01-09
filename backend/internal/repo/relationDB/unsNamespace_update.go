package relationDB

import (
	"backend/internal/types"
	"backend/share/base"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (p UnsNamespaceRepo) UpdateModelFieldsById(db *gorm.DB, id int64, fields []types.FieldDefine, numberCount int, updateAt time.Time) (int64, error) {
	db = p.model(db)
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return 0, fmt.Errorf("marshal fields error: %v", err)
	}
	result := db.
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"fields":        fieldsJSON,
			"number_fields": numberCount,
			"update_at":     updateAt,
		})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
func (p UnsNamespaceRepo) UpdateDescByAlia(db *gorm.DB, alias string, description string, updateAt time.Time) (int64, error) {
	db = p.model(db)
	result := db.
		Where("alias = ?", alias).
		Updates(map[string]interface{}{
			"description": description,
			"update_at":   updateAt,
		})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
func (p UnsNamespaceRepo) UpdateNullTemplateIdByIds(db *gorm.DB, ids []int64) (int64, error) {
	db = p.model(db)
	result := db.Where("id in ?", ids).UpdateColumn("model_id", nil)
	return result.RowsAffected, nil
}
func (p UnsNamespaceRepo) LinkLabelOnUns(db *gorm.DB, unsID, labelID int64, labelName string, updateAt time.Time) error {
	// 使用原生 SQL 片段
	sql := fmt.Sprintf(
		` jsonb_set(CASE WHEN label_ids IS NULL THEN '{}' ELSE label_ids END, '{"%d"}', '"%s"')`,
		labelID, labelName,
	)
	result := p.model(db).
		Where("id = ?", unsID).
		Updates(map[string]interface{}{
			"label_ids": gorm.Expr(sql),
			"update_at": updateAt,
		})

	return result.Error
}
func (p UnsNamespaceRepo) UnlinkLabelsByIds(db *gorm.DB, labelId int64, unsIds []int64, updateAt time.Time) (int64, error) {
	db = p.model(db)
	jsonRemoveOp := gorm.Expr(fmt.Sprintf("label_ids - '%d' ", labelId))
	// 执行更新操作
	result := db.
		Where("id IN (?) AND label_ids IS NOT NULL", unsIds).
		Updates(map[string]interface{}{
			"label_ids": jsonRemoveOp,
			"update_at": updateAt,
		})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateNamespaceLabel 更新命名空间的标签信息
func (p UnsNamespaceRepo) UpdateNamespaceLabel(db *gorm.DB, id int64, labelId string, labelName string, updateAt time.Time) error {
	// 构建JSON设置操作
	jsonSetOp := gorm.Expr(
		"jsonb_set(CASE WHEN label_ids IS NULL THEN '{}'::jsonb ELSE label_ids END, ?, ?)",
		"{"+labelId+"}",
		"\""+labelName+"\"",
	)
	// 执行更新操作
	result := p.model(db).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"label_ids": jsonSetOp,
			"update_at": updateAt,
		})

	return result.Error
}
func (p UnsNamespaceRepo) UpdateUnsLabelNames(db *gorm.DB, labelID int64, labelName string) error {
	// 构建 JSONB 路径和值
	jsonPath := fmt.Sprintf(`{"%d"}`, labelID)
	jsonValue := fmt.Sprintf(`"%s"`, escapeSQL(labelName))

	// 使用原生 SQL 表达式
	result := p.model(db).
		Where(fmt.Sprintf("jsonb_exists(label_ids ,'%d') ", labelID)). // PostgreSQL JSONB 存在操作符 ??
		Update("label_ids",
			gorm.Expr("jsonb_set(label_ids, ?, ?)", jsonPath, jsonValue),
		)

	return result.Error
}

// UpdateLabelIfExists 更新指定labelId的标签名，仅当该labelId存在时
func (p UnsNamespaceRepo) UpdateLabelIfExists(db *gorm.DB, labelId string, labelName string) error {
	// 构建JSON设置操作
	jsonSetOp := gorm.Expr(
		"jsonb_set(label_ids, ?, ?)",
		"{"+labelId+"}",
		"\""+labelName+"\"",
	)

	// 执行条件更新
	result := p.model(db).
		Where("label_ids ?? ?", labelId). // PostgreSQL的??操作符检查JSON键是否存在
		Update("label_ids", jsonSetOp)

	return result.Error
}
func (p UnsNamespaceRepo) UpdateRefUns(db *gorm.DB, id int64, idDataTypes map[int64]int, updateAt time.Time) error {
	sql := &base.StringBuilder{}
	sql.Grow(128 + len(idDataTypes)*20)
	//s.Append(fmt.Sprintf(" AND a.update_at >= '%s'::timestamp", *updateStartTime))
	sql.Append("UPDATE ").Append(TableNameUnsNamespace).
		Append(" SET update_at=").Append(fmt.Sprintf("'%s'::timestamp", updateAt)).Append(", ref_uns = ")
	for i := len(idDataTypes); i > 0; i-- {
		sql.Append("jsonb_set(")
	}
	sql.Append("case when ref_uns is null then '{}' else ref_uns end")
	for unsId, dataType := range idDataTypes {
		sql.Append(",'{\"").Long(unsId).Append("\"}','").Int(dataType).Append("')")
	}
	sql.Append(" where id=").Long(id)
	return p.model(db).Raw(sql.String()).Error
}
func (p UnsNamespaceRepo) RemoveRefUns(db *gorm.DB, id int64, calcIds []int64, updateAt time.Time) error {
	sql := &base.StringBuilder{}
	sql.Grow(128 + len(calcIds)*16)
	//s.Append(fmt.Sprintf(" AND a.update_at >= '%s'::timestamp", *updateStartTime))
	sql.Append("UPDATE ").Append(TableNameUnsNamespace).
		Append(" SET update_at=").Append(fmt.Sprintf("'%s'::timestamp", updateAt)).Append(", ref_uns = ")
	for i := len(calcIds); i > 0; i-- {
		sql.Append("jsonb_set_lax(")
	}
	sql.Append("ref_uns")
	for _, calcId := range calcIds {
		sql.Append(",'{\"").Long(calcId).Append("\"}',null,true,'delete_key')")
	}
	sql.Append(" where id=").Long(id)
	return p.model(db).Raw(sql.String()).Error
}

// getValidColumns 获取结构体有效的数据库列名
func getValidColumns(model interface{}) map[string]string {
	result := make(map[string]string)
	t := reflect.TypeOf(model)

	// 如果传入的是指针，获取其指向的类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 确保是结构体类型
	if t.Kind() != reflect.Struct {
		return result
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 获取gorm标签
		gormTag := field.Tag.Get("gorm")
		if gormTag == "" {
			continue
		}

		// 检查是否被忽略的字段
		if strings.Contains(gormTag, "-") {
			continue
		}

		// 解析column名称
		columnName := parseColumnName(gormTag)
		if columnName != "" {
			result[columnName] = field.Name
		}
	}

	return result
}

// parseColumnName 解析gorm标签中的column名称
func parseColumnName(gormTag string) string {
	// 如果标签包含"-"，表示忽略该字段
	if strings.Contains(gormTag, "-") {
		return ""
	}

	// 查找column:xxx的模式
	tagParts := strings.Split(gormTag, ";")
	for _, part := range tagParts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "column:") {
			col := strings.TrimPrefix(part, "column:")
			fen := strings.Index(col, ";")
			if fen > 0 {
				col = col[:fen]
			}
			return col
		}
	}

	// 如果没有明确指定column，使用默认的蛇形命名
	return "" //toSnakeCase(fieldName)
}

var logDeleteUns = map[string]interface{}{"status": 0}

func init() {
	excludeFields := map[string]bool{
		"id":        true,
		"alias":     true,
		"name":      true,
		"path_type": true,
		"path":      true,
		"lay_rec":   true,
		"status":    true,
	}
	columnFields := getValidColumns(&UnsNamespace{})
	for ef := range excludeFields {
		delete(columnFields, ef)
	}
	for column := range columnFields {
		logDeleteUns[column] = nil
	}
}

func (p UnsNamespaceRepo) LogicDeleteByLayRec(db *gorm.DB, layRec string) error {
	return p.model(db).Where("lay_rec like '" + escapeSQL(layRec) + "%'").Updates(logDeleteUns).Error
}

func (p UnsNamespaceRepo) LogicDeleteByIds(db *gorm.DB, ids []int64) error {
	return p.model(db).Where("id IN ?", ids).Updates(logDeleteUns).Error
}
func (p UnsNamespaceRepo) LogicDeleteById(db *gorm.DB, id int64) error {
	return p.model(db).Where("id=?", id).Updates(logDeleteUns).Error
}
