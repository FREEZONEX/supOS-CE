package relationDB

import (
	"backend/share/base"
	"context"
	"io"
	"reflect"

	"github.com/zeromicro/go-zero/core/logx"
)

type NodeRedFlowUnsExporter struct {
}

func (m NodeRedFlowUnsExporter) ExportByFlowIds(ctx context.Context, flowIds []int64, w io.Writer) (err error) {
	sql := base.StringBuilder{}
	sql.Grow(100 + len(flowIds)*20)
	sql.Append(`select parent_id ,alias from supos_node_flow_models`)
	if len(flowIds) > 0 {
		sql.Append(` where parent_id IN (`)
		for _, flowId := range flowIds {
			sql.Long(flowId).Append(`,`)
		}
		sql.SetLast(')')
	}
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	SQL := sql.String()
	err = conn.CopyTo(ctx, w, "COPY ("+SQL+" ) TO STDOUT WITH CSV HEADER")
	logx.WithContext(ctx).Info("导出 FlowUns 关联数据 SQL: ", SQL, " , err=", err)
	return
}

var nodeUnsFieldIndexMap = parseTagFields(&NoderedFlowNode{}, "gorm")

func (m NodeRedFlowUnsExporter) Csv2Model(headers, vs []string) *NoderedFlowNode {
	po := &NoderedFlowNode{}
	var values = reflect.ValueOf(po).Elem()
	for i, h := range headers {
		value := vs[i]
		index, contains := nodeUnsFieldIndexMap[h]
		if contains && len(value) > 0 {
			field := values.Field(index)
			setFieldValue(field, h, value)
		}
	}
	return po
}
