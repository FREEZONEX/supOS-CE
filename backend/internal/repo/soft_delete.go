package repo

import (
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// DeletedTime 使用 Unix 时间戳表示软删除状态，并为 GORM 注入默认查询、更新与删除规则。
type DeletedTime int64

const softDeleteActive = 0

func (DeletedTime) QueryClauses(field *schema.Field) []clause.Interface {
	return []clause.Interface{softDeleteQueryClause{Field: field}}
}

func (DeletedTime) UpdateClauses(field *schema.Field) []clause.Interface {
	return []clause.Interface{softDeleteUpdateClause{Field: field}}
}

func (DeletedTime) DeleteClauses(field *schema.Field) []clause.Interface {
	settings := schema.ParseTagSetting(field.TagSettings["SOFTDELETE"], ",")
	deleteClause := softDeleteDeleteClause{
		Field:    field,
		Flag:     settings["FLAG"] != "",
		TimeType: softDeleteTimeType(settings),
	}
	if value := settings["DELETEDATFIELD"]; value != "" {
		deleteClause.DeleteAtField = field.Schema.LookUpField(value)
	}
	return []clause.Interface{deleteClause}
}

type softDeleteClause struct{}

func (softDeleteClause) Name() string {
	return ""
}

func (softDeleteClause) Build(clause.Builder) {}

func (softDeleteClause) MergeClause(*clause.Clause) {}

type softDeleteQueryClause struct {
	softDeleteClause
	Field *schema.Field
}

func (sd softDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; ok || stmt.Statement.Unscoped {
		return
	}

	if current, ok := stmt.Clauses["WHERE"]; ok {
		if where, ok := current.Expression.(clause.Where); ok && len(where.Exprs) >= 1 {
			for _, expr := range where.Exprs {
				if orConditions, ok := expr.(clause.OrConditions); ok && len(orConditions.Exprs) == 1 {
					where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
					current.Expression = where
					stmt.Clauses["WHERE"] = current
					break
				}
			}
		}
	}

	value := any(softDeleteActive)
	if sd.Field.DefaultValue == "null" {
		value = nil
	}
	stmt.AddClause(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: value},
	}})
	stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
}

type softDeleteUpdateClause struct {
	softDeleteClause
	Field *schema.Field
}

func (sd softDeleteUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		softDeleteQueryClause{Field: sd.Field}.ModifyStatement(stmt)
	}
}

type softDeleteDeleteClause struct {
	softDeleteClause
	Field         *schema.Field
	Flag          bool
	TimeType      schema.TimeType
	DeleteAtField *schema.Field
}

func (sd softDeleteDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() != 0 || stmt.Statement.Unscoped {
		return
	}

	currentTime := stmt.DB.NowFunc()
	assignments := clause.Set{}
	if stmt.Schema != nil {
		if deletedByField := stmt.Schema.LookUpField("DeletedBy"); deletedByField != nil {
			if actorID := actorIDFromDB(stmt.DB); actorID != 0 {
				assignments = append(assignments, clause.Assignment{
					Column: clause.Column{Name: deletedByField.DBName},
					Value:  actorID,
				})
				stmt.SetColumn(deletedByField.DBName, actorID, true)
			}
		}
	}
	if sd.DeleteAtField != nil {
		value := any(sd.timeToUnix(currentTime))
		if sd.DeleteAtField.GORMDataType == "time" {
			value = currentTime
		}
		assignments = append(assignments, clause.Assignment{
			Column: clause.Column{Name: sd.DeleteAtField.DBName},
			Value:  value,
		})
		stmt.SetColumn(sd.DeleteAtField.DBName, value, true)
	}

	deletedTime := sd.timeToUnix(currentTime)
	assignments = append(clause.Set{{
		Column: clause.Column{Name: sd.Field.DBName},
		Value:  deletedTime,
	}}, assignments...)
	stmt.AddClause(assignments)
	stmt.SetColumn(sd.Field.DBName, deletedTime, true)

	if stmt.Schema != nil {
		_, identityValues := schema.GetIdentityFieldValuesMap(stmt.Context, stmt.ReflectValue, stmt.Schema.PrimaryFields)
		column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, identityValues)
		if len(values) > 0 {
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: column, Values: values},
			}})
		}

		if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
			_, identityValues = schema.GetIdentityFieldValuesMap(stmt.Context, reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
			column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, identityValues)
			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{
					clause.IN{Column: column, Values: values},
				}})
			}
		}
	}

	softDeleteQueryClause{Field: sd.Field}.ModifyStatement(stmt)
	stmt.AddClauseIfNotExists(clause.Update{})
	stmt.Build(stmt.DB.Callback().Update().Clauses...)
}

func (sd softDeleteDeleteClause) timeToUnix(value time.Time) int64 {
	switch sd.TimeType {
	case schema.UnixNanosecond:
		return value.UnixNano()
	case schema.UnixMillisecond:
		return value.UnixMilli()
	default:
		return value.Unix()
	}
}

func softDeleteTimeType(settings map[string]string) schema.TimeType {
	if settings["NANO"] != "" {
		return schema.UnixNanosecond
	}
	if settings["MILLI"] != "" {
		return schema.UnixMillisecond
	}

	switch strings.ToUpper(settings["DELETEDATFIELDUNIT"]) {
	case "NANO":
		return schema.UnixNanosecond
	case "MILLI":
		return schema.UnixMillisecond
	default:
		return schema.UnixSecond
	}
}
