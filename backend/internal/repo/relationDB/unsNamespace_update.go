package relationDB

import (
	"backend/internal/types"
	"backend/share/base"
	"encoding/json"
	"fmt"
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
			"table_name":    TableNameUnsNamespace,
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
			"table_name":  TableNameUnsNamespace,
			"description": description,
			"update_at":   updateAt,
		})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
func (p UnsNamespaceRepo) UnlinkLabelsByIds(db *gorm.DB, labelId int64, unsIds []int64, updateAt time.Time) (int64, error) {
	db = p.model(db)
	jsonRemoveOp := gorm.Expr("jsonb_remove(label_ids, ?)", labelId)
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
