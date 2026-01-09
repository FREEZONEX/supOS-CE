package relationDB

import (
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/common/utils/dbpool"
	"backend/share/base"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
           (SELECT * FROM uns_namespace WHERE id in(%s) and status=1 and id>1000  and (data_type is null OR data_type<>5 ) order by lay_rec asc) 
            TO STDOUT WITH CSV HEADER`,
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
           (SELECT * FROM uns_namespace WHERE model_id in(%s) and status=1 and id>1000  and (data_type is null OR data_type<>5 ) order by id) 
            TO STDOUT WITH CSV HEADER`,
		strings.Join(base.Map(ids, func(e int64) string {
			return strconv.FormatInt(e, 10)
		}), ","))
	err = conn.CopyTo(ctx, w, query)
	return err
}
func (p UnsNamespaceRepo) ExportCsvByLayRecAndIds(ctx context.Context, layRec []string, ids []int64, w io.Writer, asc bool) error {
	query := base.StringBuilder{}
	query.Grow(64 + 64*len(layRec) + 20*len(ids))
	query.Append("COPY (SELECT * FROM ").Append(TableNameUnsNamespace).Append(" WHERE ")
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
func (p UnsNamespaceRepo) ExportTimeSeriaNoneTables(ctx context.Context, w io.Writer) error {
	query := base.StringBuilder{}
	query.Grow(128)
	query.Append("COPY (SELECT * FROM ").Append(TableNameUnsNamespace).
		Append(" WHERE path_type =2 and data_type =1 and (table_name is null or table_name='') order by id asc").
		Append(") TO STDOUT WITH CSV HEADER")
	dbPool := getDbPool()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	err = conn.CopyTo(ctx, w, query.String())
	return err
}

func (p UnsNamespaceRepo) DoExportBatch(batchSize int, exporter func(writer io.Writer), consumer func([]*UnsNamespace)) {
	if batchSize < 1 {
		return
	}
	// 创建管道
	reader, writer := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer func() {
			_ = writer.Close()
			wg.Done()
		}()
		exporter(writer)
	}()
	go func() {
		defer func() {
			_ = reader.Close()
			wg.Done()
		}()
		// 读取 CSV 表头
		csvReader := csv.NewReader(reader)
		headers, err := csvReader.Read()
		if err != nil {
			logx.Error("exportHeader failed:", err)
			return
		}
		if batchSize > 1 {
			batch := make([]*UnsNamespace, 0, batchSize)
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				unsPO := p.Csv2Model(headers, record)
				batch = append(batch, unsPO)
				if len(batch) >= batchSize {
					consumer(batch)
					batch = batch[:0]
				}
			}
			if len(batch) > 0 {
				consumer(batch)
			}
		} else if batchSize == 1 {
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				unsPO := p.Csv2Model(headers, record)
				consumer([]*UnsNamespace{unsPO})
			}
		}
	}()
	wg.Wait()
}
func (p UnsNamespaceRepo) Csv2Model(headers, vs []string) *UnsNamespace {
	po := &UnsNamespace{}
	var values reflect.Value
	for i, h := range headers {
		value := vs[i]
		if cvt, has := exportFunctions[h]; has {
			cvt(po, value)
		} else if index, contains := fieldIndexMap[h]; contains && len(value) > 0 {
			if !values.IsValid() {
				values = reflect.ValueOf(po).Elem()
			}
			field := values.Field(index)
			setFieldValue(field, h, value)
		}
	}
	return po
}

var __timeType = reflect.TypeOf(time.Time{})
var __timePType = reflect.TypeOf(&time.Time{})

func setFieldValue(field reflect.Value, fieldName, value string) {
	addr := field.Addr().Interface()
	fieldType := field.Type()
	if field.Kind() == reflect.Ptr {

		if kind := field.Type().Elem().Kind(); kind == reflect.String { // *string 类型的处理
			// 创建一个新的string指针并赋值
			// 注意：如果原指针为nil，需要分配内存；如果已有值，则直接修改指向的值
			if field.IsNil() {
				// 分配一个新的指针
				newPtr := reflect.New(fieldType.Elem())
				newPtr.Elem().SetString(value)
				field.Set(newPtr)
			} else {
				// 直接修改指针指向的值
				field.Elem().SetString(value)
			}
			return
		} else if field.Type() == __timePType {
			ts, err := datetimeutils.ParseDate(value)
			if err == nil {
				field.Set(reflect.ValueOf(&ts))
			} else {
				logx.Errorf("Set Datetime Field Err:%v, field=%s, value=%v", err, fieldName, value)
			}
			return
		} else {
			if field.IsNil() {
				newPtr := reflect.New(fieldType.Elem())
				field.Set(newPtr)
				addr = newPtr.Elem().Addr().Interface()
			} else {
				addr = field.Elem().Addr().Interface()
			}
		}
	} else if field.Kind() == reflect.String {
		field.SetString(value)
		return
	} else if fieldType == __timeType {
		ts, err := datetimeutils.ParseDate(value)
		if err == nil {
			field.Set(reflect.ValueOf(ts))
		} else {
			logx.Errorf("Set Datetime Field Err:%v, field=%s, value=%v", err, fieldName, value)
		}
		return
	}
	err := json.Unmarshal([]byte(value), addr)
	if err != nil {
		logx.Errorf("SetFieldErr:%v, field=%s, value=%v", err, fieldName, value)
	}
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
var fieldIndexMap = parseGormFields(&UnsNamespace{})

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
	exportFunctions["create_at"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 {
			var err error
			po.CreateAt, err = datetimeutils.ParseDate(val)
			if err != nil {
				logx.Errorf("set create_at Err: %v,  v=%s", err, val)
			}
		}
	}
	exportFunctions["update_at"] = func(po *UnsNamespace, val string) {
		if len(val) > 0 {
			var err error
			po.UpdateAt, err = datetimeutils.ParseDate(val)
			if err != nil {
				logx.Errorf("set update_at Err: %v,  v=%s", err, val)
			}
		}
	}
}
func parseGormFields(ts any) (fieldIndexMap map[string]int) {
	t := reflect.TypeOf(ts)
	// 如果传入的是指针，获取其指向的类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 确保是结构体类型
	if t.Kind() != reflect.Struct {
		return
	}
	fieldIndexMap = make(map[string]int, 16)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("gorm")
		if tag == "" {
			continue
		}
		name := parseColumnName(tag)
		if len(name) > 0 {
			fieldIndexMap[name] = i
		}
	}
	return
}
