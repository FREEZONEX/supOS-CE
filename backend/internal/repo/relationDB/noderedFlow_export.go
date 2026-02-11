package relationDB

import (
	"backend/internal/common/utils/PathUtil"
	"backend/share/base"
	"context"
	"io"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type NodeRedFlowExporter struct{}

func (m NodeRedFlowExporter) ExportByGroupAndIds(ctx context.Context, groupIds []int64, ids []int64, srcFlow bool, w io.Writer) (err error) {
	sql := nodeflowGetCombineSQL(groupIds, ids, srcFlow)
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	logx.WithContext(ctx).Info("导出FlowSQL: ", sql)
	sql = strings.TrimSpace(sql)
	if sql[len(sql)-1:] == ";" {
		sql = sql[:len(sql)-1]
	}
	err = conn.CopyTo(ctx, w, "COPY ("+sql+" ) TO STDOUT WITH CSV HEADER")
	return
}

var nodeflowFieldIndexMap = parseTagFields(&NoderedFlow{}, "gorm", "json")

func (m NodeRedFlowExporter) Csv2Model(headers, vs []string) *NoderedFlow {
	po := &NoderedFlow{}
	var values = reflect.ValueOf(po).Elem()
	for i, h := range headers {
		value := vs[i]
		index, contains := nodeflowFieldIndexMap[h]
		if !contains {
			snh := PathUtil.SnakeToCamel(h)
			if snh != h {
				h = snh
				index, contains = nodeflowFieldIndexMap[h]
			}
		}
		if contains && len(value) > 0 {
			field := values.Field(index)
			setFieldValue(field, h, value)
		}
	}
	return po
}

func nodeflowGetCombineSQL(groupIds []int64, ids []int64, srcFlow bool) string {
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
		sql.Append(" nodeflow_ids (id) AS ( VALUES ")
		for i, id := range ids {
			if i > 0 {
				sql.Append(",(").Long(id).Append(")")
			} else {
				sql.Append("(").Long(id).Append("::bigint)")
			}
		}
		sql.Append("),")

		sql.Append(`
         -- 所有需要查询的分组ID（包括从nodeflow中获取的）
        all_group_ids(id) AS (`)
		if hasGroup {
			sql.Append(`SELECT id FROM group_ids UNION	`)
		}
		sql.Append(`  SELECT DISTINCT group_id   FROM supos_node_flows 
          WHERE id IN (SELECT id FROM nodeflow_ids) AND group_id IS NOT NULL   ),`)
	}

	if !hasGroup && !hasIds {
		sql.Append(`WITH `)
	}
	sql.Append(`group_data AS (
    SELECT 'group' as export_type,0 as group_id,sort,
        id,NULL as flow_id, name as flow_name, NULL as flow_status, NULL as flow_data, description, NULL as creator,g.update_at as update_time, g.create_at as create_time,
        g.id as original_group_id,
        1 as display_order,
        g.id as sort_group_id
    FROM resource_group g
    WHERE `)
	sql.Append(`g."type"=`).Int(base.SanYuan(srcFlow, 1, 2))
	if hasGroup || hasIds {
		sql.Append(` AND g.id IN (SELECT id FROM all_group_ids)`)
	}
	sql.Append(`
),
-- nodeflow数据
nodeflow_data AS (
    SELECT  '' as export_type,group_id,0 as sort,
        id, flow_id, flow_name,flow_status,flow_data,description,creator,update_time,create_time,
        d.group_id as original_group_id,
        2 as display_order,
        COALESCE(d.group_id, 0) as sort_group_id
    FROM supos_node_flows d WHERE "template"='`).Append(base.SanYuan(srcFlow, "node-red", "event-flow")).Append(`'`)

	if hasGroup || hasIds {
		sql.Append(` AND (`)
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
			sql.Append(` d.id IN (SELECT id FROM nodeflow_ids) `)
		}
		sql.Append(`)`)
	}
	sql.Append(`
),
-- 合并数据
combined_data AS (
    SELECT * FROM group_data 
    UNION ALL
    SELECT * FROM nodeflow_data
)
-- 最终结果
SELECT 
    export_type, group_id,sort, id, flow_id, flow_name,flow_status,flow_data,description,creator,update_time,create_time
FROM combined_data
ORDER BY 
    sort_group_id,
    display_order,
    id;
    `)
	return sql.String()
}
