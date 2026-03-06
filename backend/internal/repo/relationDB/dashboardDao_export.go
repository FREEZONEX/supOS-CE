package relationDB

import (
	"backend/internal/common/utils/PathUtil"
	"backend/share/base"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func (m DashboardMapper) ExportByGroupAndIds(ctx context.Context, groupIds []int64, ids []string, w io.Writer) (err error) {
	sql := dashboardGetCombineSQL(groupIds, ids)
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	logx.WithContext(ctx).Info("导出面板SQL: ", sql)
	sql = strings.TrimSpace(sql)
	if sql[len(sql)-1:] == ";" {
		sql = sql[:len(sql)-1]
	}
	err = conn.CopyTo(ctx, w, "COPY ("+sql+" ) TO STDOUT WITH CSV HEADER")
	return
}

var dashboardFieldIndexMap = parseTagFields(&DashboardModel{}, "gorm", "json")

func (m DashboardMapper) Csv2Model(headers, vs []string) *DashboardModel {
	po := &DashboardModel{}
	var values = reflect.ValueOf(po).Elem()
	for i, h := range headers {
		value := vs[i]
		index, contains := dashboardFieldIndexMap[h]
		if !contains {
			snh := PathUtil.SnakeToCamel(h)
			if snh != h {
				h = snh
				index, contains = dashboardFieldIndexMap[h]
			}
		}
		if contains && len(value) > 0 {
			field := values.Field(index)
			setFieldValue(field, h, value)
		}
	}
	return po
}

func dashboardGetCombineSQL(groupIds []int64, ids []string) string {
	hasGroup := len(groupIds) > 0
	hasIds := len(ids) > 0
	sql := base.StringBuilder{}
	sql.Grow(3096)
	if hasGroup {
		if hasIds {
			sql.Append("WITH group_ids (id) AS (  VALUES ")
		} else {
			sql.Append("WITH all_group_ids(id)  AS (  VALUES ")
		}
		for i, groupId := range groupIds {
			if i > 0 {
				sql.Append(",(").Long(groupId).Append(")")
			} else {
				sql.Append("(").Long(groupId).Append("::bigint)")
			}
		}
		sql.Append("),")
	}

	if hasIds {
		if !hasGroup {
			sql.Append(`WITH `)
		}
		sql.Append(" dashboard_ids (id) AS ( VALUES ")
		for i, id := range ids {
			str := fmt.Sprintf("'%s'", strings.ReplaceAll(id, "'", "''"))
			if i > 0 {
				sql.Append(",(").Append(str).Append(")")
			} else {
				sql.Append("(").Append(str).Append("::varchar)")
			}
		}
		sql.Append("),")

		sql.Append(`
         -- 所有需要查询的分组ID（包括从dashboard中获取的）
        all_group_ids(id) AS (`)
		if hasGroup {
			sql.Append(`SELECT id FROM group_ids UNION	`)
		}
		sql.Append(`  SELECT DISTINCT group_id   FROM uns_dashboard 
          WHERE id IN (SELECT id FROM dashboard_ids) AND group_id IS NOT NULL   ),`)
	}

	if !hasGroup && !hasIds {
		sql.Append(`WITH `)
	}
	sql.Append(`group_data AS (
    SELECT 
        CAST(g.id AS VARCHAR) as id,
        g.name,
        0 as type,
        'group' as export_type,
        sort,
        0 as group_id,
        false as need_init,
        g.description,
        NULL as creator,
        g.update_at as update_time,
        g.create_at as create_time,
        g.id as original_group_id,
        1 as display_order,
        g.id as sort_group_id
    FROM resource_group g
    WHERE g."type"=3 `)
	if hasGroup || hasIds {
		sql.Append(`AND g.id IN (SELECT id FROM all_group_ids)`)
	}
	sql.Append(`
),
-- dashboard数据
dashboard_data AS (
    SELECT 
        d.id,
        d.name,
        d.type,
        ''  as export_type,
        0 as sort,
        d.group_id,
        d.need_init,
        d.description,
        d.creator,
        d.update_time,
        d.create_time,
        d.group_id as original_group_id,
        2 as display_order,
        COALESCE(d.group_id, 0) as sort_group_id
    FROM uns_dashboard d `)

	if hasGroup || hasIds {
		sql.Append(` WHERE (`)
		if hasGroup {
			if hasIds {
				sql.Append(` d.group_id IN (SELECT id FROM group_ids)`)
			} else {
				sql.Append(` d.group_id IN (SELECT id FROM all_group_ids)`)
			}
		}
		if hasGroup && hasIds {
			sql.Append(` OR `)
		}
		if hasIds {
			sql.Append(` d.id IN (SELECT id FROM dashboard_ids) `)
		}
		sql.Append(`)`)
	}
	sql.Append(`
),
-- 合并数据
combined_data AS (
    SELECT * FROM group_data 
    UNION ALL
    SELECT * FROM dashboard_data
)
-- 最终结果
SELECT 
    export_type, group_id, id, name, type,  need_init, description, 
    creator, update_time, create_time
FROM combined_data
ORDER BY 
    sort_group_id,
    display_order,
    id;
    `)
	return sql.String()
}
