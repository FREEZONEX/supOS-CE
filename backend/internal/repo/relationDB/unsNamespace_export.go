package relationDB

import (
	"backend/internal/common/utils/dbpool"
	"backend/share/base"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

var selectColumns = []string{"id", "path_type", "parent_id", "model_id", "alias", "name", "display_name",
	"expression", "description", "label_ids",
	"protocol", "refers", "data_type", "parent_data_type", "with_flags", "fields"}

func (p UnsNamespaceRepo) ExportCsv(ctx context.Context, pathTypes []int16, w io.Writer) error {
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf(`COPY 
           (SELECT %s FROM uns_namespace WHERE path_type in(%s) and status=1 and id>1000  and (data_type is null OR data_type<>5 ) order by lay_rec asc) 
            TO STDOUT WITH CSV HEADER`, strings.Join(selectColumns, ","),
		strings.Join(base.Map(pathTypes, func(e int16) string {
			return strconv.Itoa(int(e))
		}), ","))
	err = conn.CopyTo(ctx, w, query)
	return err
}
func (p UnsNamespaceRepo) ExportCsvByIds(ctx context.Context, ids []int64, w io.Writer) error {
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf(`COPY 
           (SELECT %s FROM uns_namespace WHERE id in(%s) and status=1 and id>1000  and (data_type is null OR data_type<>5 ) order by lay_rec asc) 
            TO STDOUT WITH CSV HEADER`, strings.Join(selectColumns, ","),
		strings.Join(base.Map(ids, func(e int64) string {
			return strconv.FormatInt(e, 10)
		}), ","))
	err = conn.CopyTo(ctx, w, query)
	return err
}
func (p UnsNamespaceRepo) ExportCsvByTemplateIds(ctx context.Context, ids []int64, w io.Writer) error {
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf(`COPY 
           (SELECT %s FROM uns_namespace WHERE model_id in(%s) and status=1 and id>1000  and (data_type is null OR data_type<>5 ) order by id) 
            TO STDOUT WITH CSV HEADER`, strings.Join(selectColumns, ","),
		strings.Join(base.Map(ids, func(e int64) string {
			return strconv.FormatInt(e, 10)
		}), ","))
	err = conn.CopyTo(ctx, w, query)
	return err
}
func (p UnsNamespaceRepo) ExportCsvByLayRecAndIds(ctx context.Context, layRec []string, ids []int64, w io.Writer, asc bool) error {
	query := base.StringBuilder{}
	query.Grow(64 + 64*len(layRec) + 20*len(ids))
	query.Append("COPY (SELECT ")
	for i, col := range selectColumns {
		if i > 0 {
			query.Append(", ")
		}
		query.Append(col)
	}
	query.Append(" FROM ").Append(TableNameUnsNamespace).Append(" WHERE ")
	if len(layRec) > 0 || len(ids) > 0 {
		query.Append(" ( ")
		if len(layRec) > 0 {
			query.Append(" ( ")
			for i, lay := range layRec {
				if i > 0 {
					query.Append(" OR ")
				}
				query.Append("lay_rec like '").Append(lay).Append("%'")
			}
			query.Append(")")
		}
		if len(ids) > 0 {
			if len(layRec) > 0 {
				query.Append(" OR ")
			}
			query.Append(" id IN (")
			for _, id := range ids {
				query.Long(id).Append(",")
			}
			query.SetLast(')')
		}
		query.Append(" ) AND ")
	}
	query.Append(" status=1 and id>1000  and (data_type is null OR data_type<>5 ) order by lay_rec ").Append(base.SanYuan(asc, " ASC )", " DESC )")).
		Append(" TO STDOUT WITH CSV HEADER")

	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	err = conn.CopyTo(ctx, w, query.String())
	return err
}
func (p UnsNamespaceRepo) Csv2Model(headers, vs []string) *UnsNamespace {
	po := &UnsNamespace{}
	for i, h := range headers {
		if cvt, has := exportFunctions[h]; has {
			cvt(po, vs[i])
		}
	}
	return po
}

var _dbPool *dbpool.ExporterPool
var dbPoolOnce sync.Once

func getDbPool() *dbpool.ExporterPool {
	if _dbPool == nil {
		dbPoolOnce.Do(func() {
			pool, er := dbpool.NewExporterPool(context.Background(), dbConfig.DSN, "uns")
			if er != nil {
				log.Panicln("pg init Err", er)
			} else {
				_dbPool = pool
			}
		})
	}
	return _dbPool
}

var exportFunctions = make(map[string]func(po *UnsNamespace, val string), 32)

func init() {
	exportFunctions["id"] = func(po *UnsNamespace, val string) {
		id, _ := strconv.ParseInt(val, 10, 64)
		po.Id = id
	}
	exportFunctions["path_type"] = func(po *UnsNamespace, val string) {
		Int, er := strconv.Atoi(val)
		if er == nil {
			po.PathType = int16(Int)
		}
	}
	exportFunctions["parent_id"] = func(po *UnsNamespace, val string) {
		pid, er := strconv.ParseInt(val, 10, 64)
		if er == nil {
			po.ParentId = &pid
		}
	}
	exportFunctions["model_id"] = func(po *UnsNamespace, val string) {
		mid, er := strconv.ParseInt(val, 10, 64)
		if er == nil {
			po.ModelId = &mid
		}
	}
	exportFunctions["alias"] = func(po *UnsNamespace, val string) {
		po.Alias = val
	}
	exportFunctions["name"] = func(po *UnsNamespace, val string) {
		po.Name = val
	}
	exportFunctions["display_name"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 {
			po.DisplayName = &val
		}
	}
	exportFunctions["expression"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 {
			po.Expression = &val
		}
	}
	exportFunctions["description"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 {
			po.Description = &val
		}
	}
	exportFunctions["label_ids"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 && val[0] == '{' {
			_ = json.Unmarshal([]byte(val), &po.LabelIds)
		}
	}
	exportFunctions["protocol"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 && val[0] == '{' {
			po.Protocol = &val
		}
	}
	exportFunctions["refers"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 && val[0] == '[' {
			_ = json.Unmarshal([]byte(val), &po.Refers)
		}
	}
	exportFunctions["data_type"] = func(po *UnsNamespace, val string) {
		dt, er := strconv.Atoi(val)
		if er == nil {
			po.DataType = base.V2p(int16(dt))
		}
	}
	exportFunctions["parent_data_type"] = func(po *UnsNamespace, val string) {
		dt, er := strconv.Atoi(val)
		if er == nil {
			po.ParentDataType = base.V2p(int16(dt))
		}
	}
	exportFunctions["with_flags"] = func(po *UnsNamespace, val string) {
		dt, er := strconv.Atoi(val)
		if er == nil {
			po.WithFlags = base.V2p(int32(dt))
		}
	}
	exportFunctions["fields"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 {
			_ = json.Unmarshal([]byte(val), &po.Fields)
		}
	}
}
