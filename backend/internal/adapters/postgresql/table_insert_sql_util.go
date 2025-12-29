package postgresql

import (
	"backend/internal/types"
	"backend/share/base"
	"fmt"
	"strconv"
	"time"
)

func getInsertStatement(uns *types.CreateTopicDto, data []map[string]string) (sql string, params []interface{}) {

	sw := base.StringBuilder{}
	sw.Grow(64 + len(uns.Fields)*32)

	table := getFullTableName(uns.GetTable())
	sw.Append("INSERT INTO ").Append(table).Append(" AS t(")

	columns := uns.Fields
	for _, fd := range columns {
		sw.Append(`"`).Append(fd.Name).Append(`",`)
	}
	sw.SetLast(')').Append(" VALUES ")

	data = DeduplicationById(uns, data)

	params = make([]interface{}, 0, len(data))

	i := 0
	for _, record := range data {
		sw.Append("(")
		for _, fd := range columns {
			if v, has := record[fd.Name]; has {
				i++
				sw.Append(fmt.Sprintf("$%d", i))
				if fd.Type == types.FieldTypeDatetime {
					mill, _ := strconv.ParseFloat(v, 64)
					if mill > 0 {
						utcTime := time.UnixMilli(int64(mill)).UTC()
						v = utcTime.Format(time.RFC3339)
					}
				}
				params = append(params, v)
			} else {
				if fd.GetType() == types.FieldTypeDatetime {
					v = "NOW()"
				} else {
					v = "DEFAULT"
				}
				sw.Append(v)
			}
			sw.Append(",")
		}
		sw.SetLast(')').Append(",")
	}
	sw.SetLast(' ')
	pks := uns.GetPrimaryField()
	if len(pks) > 0 {
		sw.Append(` ON CONFLICT (`)
		for i, f := range pks {
			if i > 0 {
				sw.Append(`, `)
			}
			sw.Append(`"`).Append(f).Append(`"`)
		}
		sw.Append(`)`)
		if len(uns.Fields) > len(pks) {
			sw.Append(" DO UPDATE SET ")
			GetUpdateColumns(uns, &sw)
		} else {
			sw.Append(" DO NOTHING ")
		}
	}
	return sw.String(), params
}
func GetUpdateColumns(uns *types.CreateTopicDto, updateColumns *base.StringBuilder) {
	// 构建插入字段 和 更新字段
	updateColumns.Grow(len(uns.Fields) * 32)
	firstUpdate := false
	for _, f := range uns.Fields {
		fieldName := f.Name
		if !f.IsUnique() {
			if firstUpdate {
				updateColumns.Append(`, `)
			} else {
				firstUpdate = true
			}
			updateColumns.Append(`"`).Append(fieldName).Append(`"  = COALESCE(EXCLUDED."`).
				Append(fieldName).Append(`", t."`).Append(fieldName).Append(`")`)
		}
	}
}
