package relationDB

import (
	"backend/internal/common/utils/PathUtil"
	"backend/share/base"
	"context"
	"io"
	"reflect"

	"github.com/zeromicro/go-zero/core/logx"
)

func (m DashboardRefMapper) ExportByGroupAndIds(ctx context.Context, groupIds []int64, ids []string, w io.Writer) (err error) {
	SQL := base.StringBuilder{}
	SQL.Grow(128)
	SQL.Append("select * from uns_dashboard_ref udr")
	if hasGrp, hasIds := len(groupIds) > 0, len(ids) > 0; hasGrp || hasIds {
		SQL.Append(" join (select id from uns_dashboard where ")
		if hasIds {
			SQL.Append(" id in (")
			for _, id := range ids {
				SQL.Append(`'`).Append(id).Append(`',`)
			}
			SQL.SetLast(')')
		}
		if hasGrp {
			if hasIds {
				SQL.Append(" AND ")
			}
			SQL.Append(" group_id IN (")
			for _, groupId := range groupIds {
				SQL.Long(groupId).Append(`,`)
			}
			SQL.SetLast(')')
		}
		SQL.Append(" )ud on udr.dashboard_id =ud.id")
	}
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	sql := SQL.String()
	logx.WithContext(ctx).Info("导出dashRefSQL: ", sql)
	if sql[len(sql)-1:] == ";" {
		sql = sql[:len(sql)-1]
	}
	err = conn.CopyTo(ctx, w, "COPY ("+sql+" ) TO STDOUT WITH CSV HEADER")
	return
}

var dashFieldIndexMap = parseTagFields(&DashboardRefModel{}, "gorm", "json")

func (m DashboardRefMapper) Csv2Model(headers, vs []string) *DashboardRefModel {
	po := &DashboardRefModel{}
	var values = reflect.ValueOf(po).Elem()
	for i, h := range headers {
		value := vs[i]
		index, contains := dashFieldIndexMap[h]
		if !contains {
			snh := PathUtil.SnakeToCamel(h)
			if snh != h {
				h = snh
				index, contains = dashFieldIndexMap[h]
			}
		}
		if contains && len(value) > 0 {
			field := values.Field(index)
			setFieldValue(field, h, value)
		}
	}
	return po
}
