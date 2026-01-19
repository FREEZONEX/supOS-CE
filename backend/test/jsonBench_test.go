package test

import (
	"backend/internal/common/utils/datetimeutils"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const SYS_FIELD_CREATE_TIME = "timeStamp"

func TestJsonGetTs(t *testing.T) {
	jsonStr := `{"double1":-20.791168813592314,"quality":0,"timeStamp":1768285507000 }`
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		t.Fatal(err)
	}
	ts := time.Now()
	m[SYS_FIELD_CREATE_TIME] = &ts
	t.Log("ts:", parseTimestamp(m[SYS_FIELD_CREATE_TIME]))
}
func newErr() (er error) {
	return er
}
func parseTimestamp(curT any) (ct int64) {
	if curT == nil {
		return -1
	} else if Float, isFloat := curT.(float64); isFloat { // json unmarshal 来的都是 float64 类型
		ct = int64(Float)
	} else if Long, isLong := curT.(int64); isLong {
		ct = Long
	} else {
		str := fmt.Sprint(curT)
		Double, err := strconv.ParseFloat(str, 64)
		if err != nil {
			ct = -1
			if dt, dtEr := datetimeutils.ParseDate(str); dtEr == nil && dt.Year() > 1970 {
				ct = dt.UnixMilli()
			}
		} else if Int := int64(Double); Int > 1100000000000 {
			ct = Int
		}
	}
	if ct < 1100000000000 || ct > 11000000000001 {
		return -1
	}
	return ct
}
func BenchmarkJson(b *testing.B) {
	jsonStr := `{"timeStamp":1768288504446, "wet":168.931,"qos":0}`
	b.StartTimer()
	//BenchmarkJson-12    	  683923	      1636 ns/op
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &m)
		if err != nil {
			b.Fatal(err)
		}
		ts, hasTs := m[SYS_FIELD_CREATE_TIME]
		value, hasValue := m["wet"].(float64)
		if !hasTs || !hasValue {
			b.Fatal("wet", ts, value)
		}
	}
}
func Benchmark_jsoniter(b *testing.B) {
	jsonStr := `{"timeStamp":1768288504446, "wet":168.931,"qos":0}`
	b.StartTimer()
	//Benchmark_jsoniter-12    	 1504185	       783.6 ns/op
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &m)
		if err != nil {
			b.Fatal(err)
		}
		ts, hasTs := m[SYS_FIELD_CREATE_TIME]
		value, hasValue := m["wet"].(float64)
		if !hasTs || !hasValue {
			b.Fatal("wet", ts, value)
		}
	}
}
