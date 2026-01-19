package relationDB

import (
	"backend/internal/common/utils/datetimeutils"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/conf"
)

func parsePostgresTimestamp(tc string) (time.Time, error) {
	plus := strings.LastIndex(tc, "+")
	zone := 0
	hasZone := false
	if plus > 0 {
		var err error
		zone, err = strconv.Atoi(tc[plus+1:])
		if err == nil {
			hasZone = true
		}
		tc = strings.TrimSpace(tc[:plus])
	}
	tm, err := time.Parse("2006-01-02 15:04:05.999999", tc)
	if err == nil {
		if hasZone {
			local := time.FixedZone("", zone)
			tm = tm.In(local)
		}
		return tm, nil
	}
	return time.Time{}, err
}

func Test_parsePgData(t *testing.T) {
	// 测试不同的 PostgreSQL 时间格式
	testCases := []string{
		"2026-01-08 16:46:18.523 +0800",
		"2026-01-06 13:04:16.379395+00",
		"2026-01-06 22:04:16.379395+08",
	}
	{
		tc := "2026-01-08 16:46:18.523 +0800"
		ts, err := time.Parse("2006-01-02 15:04:05.999 Z0700", tc)
		if err != nil {
			t.Logf("解析失败 %s: %v\n", tc, err)
		} else {
			t.Logf("解析成功 %s -> %v (UTC: %v)\n", tc, ts, ts.UTC())
		}
	}
	for _, tc := range testCases {
		ts, err := datetimeutils.ParseDate(tc)
		if err != nil {
			t.Logf("解析失败 %s: %v\n", tc, err)
		} else {
			t.Logf("解析成功 %s -> %v (UTC: %v)\n", tc, ts, ts.UTC())
		}
	}

}
func Test_setFieldValue(t *testing.T) {
	var po = &UnsNamespace{}

	headers := []string{"Id", "Name", "ParentId", "DataSrcId", "CreateAt"}
	values := []string{"2006302406467391607", "opcua_demo_312", "2008499350988132352", "3", "2026-01-08 16:46:18.523 +0800"}
	refValues := reflect.ValueOf(po).Elem()
	for i, header := range headers {
		value := values[i]
		field := refValues.FieldByName(header)
		setFieldValue(field, header, value)
	}
	bs, _ := json.MarshalIndent(po, "", " ")
	t.Log(string(bs))
}
func Test_Csv2Model(t *testing.T) {
	InitDbConfig(conf.Database{
		DBType: "pgsql",
		DSN:    "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos",
	})

	var mapper UnsNamespaceRepo
	{
		headers := []string{"id", "name", "path_type", "data_src_id", "create_at"}
		values := []string{"2006302406467391607", "opcua_demo_312", "2", "3", "2026-01-08 16:46:18.523 +0800"}
		po := mapper.Csv2Model(headers, values)
		bs, _ := json.MarshalIndent(po, "", " ")
		t.Log(string(bs))
	}

	mapper.DoExportBatch(100, func(writer io.Writer) {
		mapper.ExportCsvByIds(context.Background(), []int64{2008499350988132353, 2008743696706572381}, writer)
	}, func(namespaces []*UnsNamespace) {
		bs, _ := json.MarshalIndent(namespaces, "", " ")
		t.Log(string(bs))
	})
}
