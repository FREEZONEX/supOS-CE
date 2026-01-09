package timescaledb

import (
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"testing"
	"time"
)

// 测试辅助函数
func createField(name, fieldType string, index *string) *types.FieldDefine {
	return &types.FieldDefine{
		Name:  name,
		Type:  fieldType,
		Index: index,
	}
}

func strPtr(s string) *string {
	return &s
}

// 为了测试，我们需要包装成 serviceApi.UnsData 的格式
// 假设 serviceApi.UnsData 结构如下：
// type UnsData struct {
//     Uns  UnsInfo
//     Data []map[string]string
// }

// TestPreprocess_EmptyData 测试空数据
func TestPreprocess_EmptyData(t *testing.T) {
	var unsData []serviceApi.UnsData

	result := preprocess(unsData)

	if result.conflict.columns != nil || result.conflict.rows != nil {
		t.Errorf("Expected empty conflict result")
	}
	if result.normal.columns != nil || result.normal.rows != nil {
		t.Errorf("Expected empty normal result")
	}
}
func procFields(fs []*types.FieldDefine) []*types.FieldDefine {
	name := constants.SystemSeqTag // Ensure the name is correct
	tableValueField := &types.FieldDefine{
		Name:        constants.SystemSeqTag,
		Type:        types.FieldTypeLong,
		Unique:      base.OptionalTrue,
		TbValueName: &name,
	}
	fs = append(fs,
		&types.FieldDefine{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime, Unique: base.OptionalTrue},
		tableValueField,
		&types.FieldDefine{Name: constants.QosField, Type: types.FieldTypeLong},
	)
	return fs
}

// TestPreprocess_SingleUnsNoConflict 测试初始单个 Uns,默认会有冲突
func TestPreprocess_SingleUnsNoConflict(t *testing.T) {
	unsInfo := &types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("temperature", types.FieldTypeFloat, base.V2p("double_1")),
			createField("humidity", types.FieldTypeFloat, base.V2p("double_2")),
			createField("state", types.FieldTypeString, base.V2p("str_1")),
		}),
		Timestamps: [2]int64{0, 0},
	}

	testData := []serviceApi.UnsData{
		{
			Uns: unsInfo,
			Data: []map[string]string{
				{
					constants.SysFieldCreateTime: "1000.5",
					"temperature":                "25.5",
					"humidity":                   "60.0",
					"state":                      "normal",
				},
				{
					constants.SysFieldCreateTime: "2000.0",
					"temperature":                "26.0",
					"humidity":                   "62.0",
					"state":                      "normal",
				},
			},
		},
	}

	result := preprocess(testData)

	// 验证 正常数据应该为空
	if result.normal.rows != nil && len(result.normal.rows) > 0 {
		t.Errorf("Expected no normal rows, got %d", len(result.normal.rows))
	}

	// 验证冲突数据应该有两行
	if result.conflict.rows == nil || len(result.conflict.rows) != 2 {
		t.Errorf("Expected 2 conflict rows, got %v", result.conflict.rows)
	}

	// 验证列名
	expectedColumnsStart := []string{
		constants.SysFieldCreateTime,
		constants.SystemSeqTag,
		constants.QosField,
	}

	// 验证系统字段在前
	for i, col := range expectedColumnsStart {
		if result.conflict.columns[i] != col {
			t.Errorf("Expected column %s at index %d, got %s", col, i, result.conflict.columns[i])
		}
	}

	// 验证自定义字段排序正确（按字母顺序）
	customCols := result.conflict.columns[3:]
	for i := 0; i < len(customCols)-1; i++ {
		if customCols[i] > customCols[i+1] {
			t.Errorf("Columns not sorted: %s > %s", customCols[i], customCols[i+1])
		}
	}
}

// TestPreprocess_WithConflictData 测试有冲突的数据
func TestPreprocess_WithConflictData(t *testing.T) {
	unsInfo := &types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("value", types.FieldTypeFloat, base.V2p("double_1")),
		}),
		Timestamps: [2]int64{0, 1000}, // 设置阈值为1000
	}

	testData := []serviceApi.UnsData{
		{
			Uns: unsInfo,
			Data: []map[string]string{
				{
					constants.SysFieldCreateTime: "500.0", // 小于阈值，应该进入conflict
					"value":                      "100.5",
				},
				{
					constants.SysFieldCreateTime: "1500.0", // 大于阈值，应该进入normal
					"value":                      "200.5",
				},
				{
					constants.SysFieldCreateTime: "500.0", // 相同时间戳，应该合并到conflict
					"value":                      "150.5",
				},
			},
		},
	}

	result := preprocess(testData)

	// 验证冲突数据有一行（合并后）
	if result.conflict.rows == nil || len(result.conflict.rows) != 1 {
		t.Errorf("Expected 1 conflict row, got %d", len(result.conflict.rows))
	}

	// 验证正常数据有一行
	if result.normal.rows == nil || len(result.normal.rows) != 1 {
		t.Errorf("Expected 1 normal row, got %d", len(result.normal.rows))
	}

	// 验证冲突数据合并（后一个值覆盖前一个）
	if result.conflict.rows[0][3] != "150.5" {
		t.Errorf("Expected conflict value 150.5, got %v", result.conflict.rows[0][3])
	}

	// 验证正常数据
	if result.normal.rows[0][3] != "200.5" {
		t.Errorf("Expected normal value 200.5, got %v", result.normal.rows[0][3])
	}
}

// TestPreprocess_DatetimeFieldConversion 测试datetime类型字段转换
func TestPreprocess_DatetimeFieldConversion(t *testing.T) {
	unsInfo := &types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("event_time", types.FieldTypeDatetime, base.V2p("date_1")),
		}),
		Timestamps: [2]int64{0, 0},
	}

	timestamp := "1672531200000" // 2023-01-01 00:00:00 UTC
	expectedTime := time.UnixMilli(1672531200000).UTC().Format("2006-01-02 15:04:05.000") + "+00"

	testData := []serviceApi.UnsData{
		{
			Uns: unsInfo,
			Data: []map[string]string{
				{
					constants.SysFieldCreateTime: "1000.0",
					"event_time":                 timestamp,
				},
			},
		},
	}

	result := preprocess(testData)

	if result.conflict.rows == nil || len(result.conflict.rows) == 0 {
		t.Fatal("No conflict rows generated")
	}

	convertedTime := result.conflict.rows[0][3] // event_time在第3列（前面3个系统字段）
	if convertedTime != expectedTime {
		t.Errorf("DateTime conversion failed. Expected %s, got %s", expectedTime, convertedTime)
	}
}

// TestPreprocess_MultipleUns 测试多个Uns的数据
func TestPreprocess_MultipleUns(t *testing.T) {
	uns1 := &types.CreateTopicDto{
		Id:        1,
		Alias:     "uns1",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("temp", types.FieldTypeFloat, base.V2p("double_1")),
		}),
		Timestamps: [2]int64{0, 0},
	}

	uns2 := &types.CreateTopicDto{
		Id:        2,
		Alias:     "uns2",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("pressure", types.FieldTypeFloat, base.V2p("double_1")),
		}),
		Timestamps: [2]int64{0, 0},
	}

	testData := []serviceApi.UnsData{
		{
			Uns: uns1,
			Data: []map[string]string{
				{constants.SysFieldCreateTime: "1000.0", "temp": "25.5"},
			},
		},
		{
			Uns: uns2,
			Data: []map[string]string{
				{constants.SysFieldCreateTime: "2000.0", "pressure": "101.3"},
			},
		},
	}

	result := preprocess(testData)

	// 计算总行数
	totalRows := 0
	if result.conflict.rows != nil {
		totalRows += len(result.conflict.rows)
	}
	if result.normal.rows != nil {
		totalRows += len(result.normal.rows)
	}

	if totalRows != 2 {
		t.Errorf("Expected total 2 rows, got %d", totalRows)
	}

	// 验证列名包含两个Uns的所有字段
	expectedColumns := []string{
		constants.SysFieldCreateTime,
		constants.SystemSeqTag,
		constants.QosField,
		"double_1",
	}

	if len(result.conflict.columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(result.normal.columns))
	}

	for i, col := range expectedColumns {
		if result.conflict.columns[i] != col {
			t.Errorf("Column %d: expected %s, got %s", i, col, result.normal.columns[i])
		}
	}
}

// TestPreprocess_TimestampThreshold 测试时间戳阈值逻辑
func TestPreprocess_TimestampThreshold(t *testing.T) {
	// 测试不同阈值情况
	testCases := []struct {
		name           string
		threshold      int64
		timestamp      string
		expectConflict bool
	}{
		{"小于阈值", 1000, "500.0", true},
		{"等于阈值", 1000, "1000.0", true},
		{"大于阈值", 1000, "1500.0", false},
		{"零阈值", 0, "500.0", false},
		{"最大阈值", 0, "500.0", false}, // 阈值为0，所有数据都进normal
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unsInfo := &types.CreateTopicDto{
				Id:        1,
				Alias:     "test_uns",
				TableName: "uns_timeserial",
				Fields: procFields([]*types.FieldDefine{
					createField("value", types.FieldTypeFloat, base.V2p("double_1")),
				}),
				Timestamps: [2]int64{0, tc.threshold},
			}

			testData := []serviceApi.UnsData{
				{
					Uns: unsInfo,
					Data: []map[string]string{
						{
							constants.SysFieldCreateTime: tc.timestamp,
							"value":                      "100.0",
						},
					},
				},
			}

			result := preprocess(testData)

			if tc.expectConflict {
				if result.conflict.rows == nil || len(result.conflict.rows) != 1 {
					t.Errorf("Expected 1 conflict row, got %d", len(result.conflict.rows))
				}
				if result.normal.rows != nil && len(result.normal.rows) > 0 {
					t.Errorf("Expected no normal rows, got %d", len(result.normal.rows))
				}
			} else {
				if result.normal.rows == nil || len(result.normal.rows) != 1 {
					t.Errorf("Expected 1 normal row, got %d", len(result.normal.rows))
				}
				if result.conflict.rows != nil && len(result.conflict.rows) > 0 {
					t.Errorf("Expected no conflict rows, got %d", len(result.conflict.rows))
				}
			}
		})
	}
}

// TestPreprocess_FieldIndexMapping 测试字段索引映射
func TestPreprocess_FieldIndexMapping(t *testing.T) {
	unsInfo := &types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("field1", types.FieldTypeString, base.V2p("str_1")),
			createField("field2", types.FieldTypeString, base.V2p("str_2")),
			createField("field3", types.FieldTypeString, base.V2p("str_3")),
		}),
		Timestamps: [2]int64{0, 0},
	}

	testData := []serviceApi.UnsData{
		{
			Uns: unsInfo,
			Data: []map[string]string{
				{
					constants.SysFieldCreateTime: "1000.0",
					"field1":                     "value1",
					"field2":                     "value2",
					"field3":                     "value3",
				},
			},
		},
	}

	result := preprocess(testData)

	// 验证列名使用了索引名
	expectedColumns := []string{
		constants.SysFieldCreateTime,
		constants.SystemSeqTag,
		constants.QosField,
		"field2", // 没有index，使用原名
		"idx1",   // 使用索引名
		"idx3",   // 使用索引名
	}

	for i, col := range expectedColumns {
		if result.normal.columns[i] != col {
			t.Errorf("Column %d: expected %s, got %s", i, col, result.normal.columns[i])
		}
	}
}

// TestPreprocess_InvalidTimestampHandling 测试无效时间戳处理
func TestPreprocess_InvalidTimestampHandling(t *testing.T) {
	unsInfo := &types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("value", types.FieldTypeFloat, base.V2p("double_1")),
		}),
		Timestamps: [2]int64{0, 0},
	}

	testData := []serviceApi.UnsData{
		{
			Uns: unsInfo,
			Data: []map[string]string{
				{constants.SysFieldCreateTime: "invalid", "value": "100"}, // 无效时间戳
				{constants.SysFieldCreateTime: "0", "value": "200"},       // 时间戳为0
				{constants.SysFieldCreateTime: "-100", "value": "300"},    // 负时间戳
				{constants.SysFieldCreateTime: "1000.5", "value": "400"},  // 有效时间戳
			},
		},
	}

	result := preprocess(testData)

	// 应该只有一行有效数据
	totalRows := 0
	if result.conflict.rows != nil {
		totalRows += len(result.conflict.rows)
	}
	if result.normal.rows != nil {
		totalRows += len(result.normal.rows)
	}

	if totalRows != 1 {
		t.Errorf("Expected 1 valid row, got %d", totalRows)
	}
}

// TestPreprocess_MaxTimestampUpdate 测试最大时间戳更新
func TestPreprocess_MaxTimestampUpdate(t *testing.T) {
	unsInfo := &types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("value", types.FieldTypeFloat, base.V2p("double_1")),
		}),
		Timestamps: [2]int64{-1, 0}, // 初始最大值为-1
	}

	testData := []serviceApi.UnsData{
		{
			Uns: unsInfo,
			Data: []map[string]string{
				{constants.SysFieldCreateTime: "500.0", "value": "100"},
				{constants.SysFieldCreateTime: "1500.0", "value": "200"},
				{constants.SysFieldCreateTime: "1000.0", "value": "300"},
			},
		},
	}

	result := preprocess(testData)

	// 验证最大时间戳更新为1500
	if unsInfo.Timestamps[0] != 1500 {
		t.Errorf("Expected max timestamp 1500, got %d", unsInfo.Timestamps[0])
	}

	// 验证所有数据都被处理
	totalRows := 0
	if result.conflict.rows != nil {
		totalRows += len(result.conflict.rows)
	}
	if result.normal.rows != nil {
		totalRows += len(result.normal.rows)
	}

	if totalRows != 3 {
		t.Errorf("Expected 3 rows, got %d", totalRows)
	}
}

// 集成测试：模拟真实场景
func TestPreprocess_Integration(t *testing.T) {
	// 创建多个UnsInfo
	uns1 := &types.CreateTopicDto{
		Id:        1,
		Alias:     "temperature_sensor",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("temperature", types.FieldTypeFloat, base.V2p("double_1")),
			createField("battery", types.FieldTypeFloat, nil),
		}),
		Timestamps: [2]int64{-1, 500}, // 阈值为500
	}

	uns2 := &types.CreateTopicDto{
		Id:        2,
		Alias:     "pressure_sensor",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("pressure", types.FieldTypeFloat, base.V2p("double_1")),
			createField("state", types.FieldTypeString, base.V2p("str_1")),
		}),
		Timestamps: [2]int64{-1, 1000}, // 阈值为1000
	}

	testData := []serviceApi.UnsData{
		{
			Uns: uns1,
			Data: []map[string]string{
				{constants.SysFieldCreateTime: "400.0", "temperature": "22.5", "battery": "95.0"},
				{constants.SysFieldCreateTime: "600.0", "temperature": "23.0", "battery": "94.5"},
				{constants.SysFieldCreateTime: "400.0", "temperature": "22.8", "battery": "95.0"}, // 冲突
			},
		},
		{
			Uns: uns2,
			Data: []map[string]string{
				{constants.SysFieldCreateTime: "900.0", "pressure": "101.3", "state": "normal"},
				{constants.SysFieldCreateTime: "1100.0", "pressure": "101.2", "state": "normal"},
			},
		},
	}

	result := preprocess(testData)

	// 验证冲突数据
	conflictRows := 0
	if result.conflict.rows != nil {
		conflictRows = len(result.conflict.rows)
	}

	// uns1 有1个冲突（400.0），uns2 有1个冲突（900.0）
	if conflictRows != 2 {
		t.Errorf("Expected 2 conflict rows, got %d", conflictRows)
	}

	// 验证正常数据
	normalRows := 0
	if result.normal.rows != nil {
		normalRows = len(result.normal.rows)
	}

	// uns1 有1个正常（600.0），uns2 有1个正常（1100.0）
	if normalRows != 2 {
		t.Errorf("Expected 2 normal rows, got %d", normalRows)
	}

	// 验证列名排序
	expectedColumns := []string{
		constants.SysFieldCreateTime,
		constants.SystemSeqTag,
		constants.QosField,
		"battery",
		"pressure",
		"st",   // status的索引名
		"temp", // temperature的索引名
	}

	if len(result.normal.columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(result.normal.columns))
	}

	for i, col := range expectedColumns {
		if result.normal.columns[i] != col {
			t.Errorf("Column %d: expected %s, got %s", i, col, result.normal.columns[i])
		}
	}
}
