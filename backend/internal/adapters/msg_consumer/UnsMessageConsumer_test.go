package msg_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestNpePanic(t *testing.T) {
	var vm map[string]interface{}
	//vm = make(map[string]interface{})
	//vm["a"] = 1
	CT := "timestamp"
	if len(vm) == 0 {
		t.Log("no vm")
	} else {
		curT, hasCt := vm[CT]
		t.Logf("curT:%v, hasCt:%v", curT, hasCt)
	}
	t.Log("ts:", parseTimestamp(vm[CT]))
}

// 模拟Java中的mockBean方法
func mockBean(p string) map[string]interface{} {
	return map[string]interface{}{
		p: true,
	}
}

// 模拟Java中的mockBean方法（带时间戳）
func mockBeanWithTimestamp(timestamp int64, p string) map[string]interface{} {
	return map[string]interface{}{
		SYS_FIELD_CREATE_TIME: timestamp,
		p:                     true,
	}
}

// 辅助函数：检查两个map的key是否相同
func keysEqual(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}

	keys1 := make([]string, 0, len(m1))
	for k := range m1 {
		keys1 = append(keys1, k)
	}

	keys2 := make([]string, 0, len(m2))
	for k := range m2 {
		keys2 = append(keys2, k)
	}

	sort.Strings(keys1)
	sort.Strings(keys2)

	return reflect.DeepEqual(keys1, keys2)
}

// 辅助函数：创建HashSet（Go中没有HashSet，用map[string]bool模拟）
func newHashSet(items ...string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

// 检查map的keys是否与预期的set匹配
func keysMatch(m map[string]interface{}, expectedKeys map[string]bool) bool {
	if len(m) != len(expectedKeys) {
		return false
	}

	for key := range m {
		if !expectedKeys[key] {
			return false
		}
	}

	return true
}

const SYS_FIELD_CREATE_TIME = "timeStamp"

func mergeBeansWithCT(list []map[string]interface{}, prevBean map[string]any) []map[string]interface{} {
	return mergeBeansWithTimestamp(context.Background(), list, SYS_FIELD_CREATE_TIME, time.Now().UnixMilli(), prevBean)
}
func TestLogPanic(t *testing.T) {
	ctx := t.Context()
	var err error
	CT := "timestamp"
	var list []map[string]interface{}
	prevBean := map[string]any{}
	prevBean[CT] = float64(time.Now().UnixMilli())
	logx.WithContext(ctx).Errorf("HandleThrow|traceID=%s|error=%#v|stack=%s| CT=%s, list=%+v, prevBean=%+v", utils.TraceIdFromContext(ctx),
		err, utils.Stack(2, 20), CT, list, prevBean)

}
func TestMergeBeansWithTimestamp2(t *testing.T) {
	floatTime := 1.766978836e+12
	ct := int64(floatTime)
	t.Log("ct=", ct)
	ctStr := strconv.FormatInt(ct, 10)
	prevJson := `{"double2":99.86295347545729, "timeStamp":` + ctStr + `}`
	curJson := `{"double1":91, "timeStamp":` + ctStr + `}`

	var prevBean map[string]interface{}
	var curBean map[string]interface{}
	json.Unmarshal([]byte(prevJson), &prevBean)
	json.Unmarshal([]byte(curJson), &curBean)
	rs := mergeBeansWithCT([]map[string]any{curBean}, prevBean)

	rsJson, _ := json.Marshal(rs)
	t.Log("rs:", string(rsJson))
}
func TestMergeBeansWithTimestamp(t *testing.T) {
	// Test case 1
	t.Run("case1_different_timestamps", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		t0 := time.Now().UnixNano() / int64(time.Millisecond)
		list = append(list, mockBeanWithTimestamp(t0, "p1"))
		list = append(list, mockBeanWithTimestamp(t0+3, "p2"))

		rs := mergeBeansWithCT(list, nil)

		if len(rs) != 2 {
			t.Errorf("Expected 2 items, got %d", len(rs))
		}

		expectedKeys1 := newHashSet(SYS_FIELD_CREATE_TIME, "p1")
		if !keysMatch(rs[0], expectedKeys1) {
			t.Errorf("Expected keys %v, got keys %v", expectedKeys1, rs[0])
		}

		expectedKeys2 := newHashSet(SYS_FIELD_CREATE_TIME, "p2")
		if !keysMatch(rs[1], expectedKeys2) {
			t.Errorf("Expected keys %v, got keys %v", expectedKeys2, rs[1])
		}
	})

	// Test case 2
	t.Run("case2_no_timestamp_in_list", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		list = append(list, mockBean("p2"))

		lastMsg := mockBean("p1")
		rs := mergeBeansWithCT(list, lastMsg)

		if len(rs) != 1 {
			t.Errorf("Expected 1 item, got %d", len(rs))
		}

		expectedKeys := newHashSet(SYS_FIELD_CREATE_TIME, "p2")
		if !keysMatch(rs[0], expectedKeys) {
			t.Errorf("Expected keys %v, got keys %v", expectedKeys, rs[0])
		}

		fmt.Printf("case 2: %v\n", list)
	})

	// Test case 3
	t.Run("case3_same_timestamp_with_lastMsg", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		t0 := time.Now().UnixNano() / int64(time.Millisecond)
		list = append(list, mockBeanWithTimestamp(t0, "p2"))

		lastMsg := mockBeanWithTimestamp(t0, "p1")
		rs := mergeBeansWithCT(list, lastMsg)

		if len(rs) != 1 {
			t.Errorf("Expected 1 item, got %d", len(rs))
		}

		jsonRs, _ := json.Marshal(rs)
		fmt.Printf("case 3: %s\n", string(jsonRs))

		expectedKeys := newHashSet(SYS_FIELD_CREATE_TIME, "p1", "p2")
		if !keysMatch(rs[0], expectedKeys) {
			t.Errorf("Expected keys %v, got keys %v", expectedKeys, rs[0])
		}
	})

	// Test case 4
	t.Run("case4_multiple_items_same_timestamp", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		t0 := time.Now().UnixNano() / int64(time.Millisecond)
		list = append(list, mockBeanWithTimestamp(t0, "p1"))
		list = append(list, mockBeanWithTimestamp(t0, "p2"))
		list = append(list, mockBeanWithTimestamp(t0, "p3"))

		rs := mergeBeansWithCT(list, nil)

		if len(rs) != 1 {
			t.Errorf("Expected 1 item, got %d", len(rs))
		}

		expectedKeys := newHashSet(SYS_FIELD_CREATE_TIME, "p1", "p2", "p3")
		if !keysMatch(rs[0], expectedKeys) {
			t.Errorf("Expected keys %v, got keys %v", expectedKeys, rs[0])
		}
	})

	// Test case 5
	t.Run("case5_multiple_items_with_lastMsg_same_timestamp", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		t0 := time.Now().UnixNano() / int64(time.Millisecond)
		list = append(list, mockBeanWithTimestamp(t0, "p2"))
		list = append(list, mockBeanWithTimestamp(t0, "p3"))

		lastMsg := mockBeanWithTimestamp(t0, "p1")
		rs := mergeBeansWithCT(list, lastMsg)

		if len(rs) != 1 {
			t.Errorf("Expected 1 item, got %d", len(rs))
		}

		expectedKeys := newHashSet(SYS_FIELD_CREATE_TIME, "p1", "p2", "p3")
		if !keysMatch(rs[0], expectedKeys) {
			t.Errorf("Expected keys %v, got keys %v", expectedKeys, rs[0])
		}
	})

	// Test case 6
	t.Run("case6_mixed_timestamps", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		t0 := time.Now().UnixNano() / int64(time.Millisecond)
		list = append(list, mockBeanWithTimestamp(t0, "p1"))
		list = append(list, mockBeanWithTimestamp(t0, "p2"))
		list = append(list, mockBeanWithTimestamp(t0+1, "p3"))

		rs := mergeBeansWithCT(list, nil)

		if len(rs) != 2 {
			t.Errorf("Expected 2 items, got %d", len(rs))
		}

		expectedKeys1 := newHashSet(SYS_FIELD_CREATE_TIME, "p1", "p2")
		if !keysMatch(rs[0], expectedKeys1) {
			t.Errorf("Expected keys %v for first item, got keys %v", expectedKeys1, rs[0])
		}

		expectedKeys2 := newHashSet(SYS_FIELD_CREATE_TIME, "p3")
		if !keysMatch(rs[1], expectedKeys2) {
			t.Errorf("Expected keys %v for second item, got keys %v", expectedKeys2, rs[1])
		}
	})

	// Test case 7
	t.Run("case7_multiple_items_mixed_timestamps", func(t *testing.T) {
		list := make([]map[string]interface{}, 0)
		t0 := time.Now().UnixNano() / int64(time.Millisecond)
		list = append(list, mockBeanWithTimestamp(t0, "p1"))
		list = append(list, mockBeanWithTimestamp(t0, "p2"))
		list = append(list, mockBeanWithTimestamp(t0+1, "p3"))
		list = append(list, mockBeanWithTimestamp(t0+1, "p4"))

		rs := mergeBeansWithCT(list, nil)

		if len(rs) != 2 {
			t.Errorf("Expected 2 items, got %d", len(rs))
		}

		expectedKeys1 := newHashSet(SYS_FIELD_CREATE_TIME, "p1", "p2")
		if !keysMatch(rs[0], expectedKeys1) {
			t.Errorf("Expected keys %v for first item, got keys %v", expectedKeys1, rs[0])
		}

		expectedKeys2 := newHashSet(SYS_FIELD_CREATE_TIME, "p3", "p4")
		if !keysMatch(rs[1], expectedKeys2) {
			t.Errorf("Expected keys %v for second item, got keys %v", expectedKeys2, rs[1])
		}
	})
}
