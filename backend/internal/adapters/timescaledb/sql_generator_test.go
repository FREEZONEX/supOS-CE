package timescaledb

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLGenerator_GenerateSyncSQLs(t *testing.T) {
	// 创建 SQL 生成器
	sqlGen := NewSQLGenerator()

	// 准备测试数据
	unsList := []UnsViewInfo{
		{
			Uns: &types.CreateTopicDto{
				Id:        1001,
				Alias:     "test_view_1",
				TableName: "uns_timeserial",
				Fields: []*types.FieldDefine{
					{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime, Unique: base.OptionalTrue},
					{Name: constants.SysFieldID, Type: types.FieldTypeLong, Unique: base.OptionalTrue},
					{Name: "temp", Type: types.FieldTypeInteger},
					{Name: "qst", Type: types.FieldTypeLong},
					{Name: constants.QosField, Type: types.FieldTypeLong},
				},
			},
			View: SimpleViewInfo{
				SrcTable: "uns_timeserial",
				Columns: []ViewColumnInfo{
					{ColumnName: "temp", SourceColumn: "double_1"},
				},
			},
		},
	}

	// 物理表字段为空（表不存在）
	var physicsTableFields []*types.FieldDefine

	// 生成 SQL
	result := sqlGen.GenerateSyncSQLs(physicsTableFields, constants.SysFieldCreateTime, constants.QosField, unsList)

	// 验证结果
	assert.False(t, result.HasErrors())
	assert.Greater(t, len(result.CreateTableSQL), 0)
	assert.Greater(t, len(result.CreateViewSQL), 0)

	// 验证字段映射
	assert.Contains(t, result.FieldMappings, "test_view_1")

	// 验证 SQL 数量
	assert.Equal(t, 5, len(result.CreateTableSQL))
	assert.Equal(t, 0, len(result.AlterTableSQL))
	assert.Equal(t, 0, len(result.UpdateDataSQL))
	assert.Equal(t, 1, len(result.CreateViewSQL))

	t.Logf("Create Table SQL:\n%s", result.CreateTableSQL[0])
	t.Logf("Create View SQL:\n%s", result.CreateViewSQL[0])
	for _, uns := range unsList {
		for _, f := range uns.Uns.GetFields() {
			t.Logf("field: %s, index=%v\n", f.Name, base.P2v(f.Index))
		}
	}
}

func TestSQLGenerator_FieldMapping(t *testing.T) {
	sqlGen := NewSQLGenerator()

	unsInfo := &types.CreateTopicDto{
		Id:        1002,
		Alias:     "test_view_2",
		TableName: "uns_timeserial",
		Fields: []*types.FieldDefine{
			{Name: "field1", Type: types.FieldTypeDouble},
			{Name: "field2", Type: types.FieldTypeLong},
			{Name: "field3", Type: types.FieldTypeDouble},
		},
	}

	viewInfo := SimpleViewInfo{
		SrcTable: "uns_timeserial",
		Columns: []ViewColumnInfo{
			{ColumnName: "field1", SourceColumn: "double_1"},
			{ColumnName: "field3", SourceColumn: "double_3"},
		},
	}

	unsView := UnsViewInfo{
		Uns:  unsInfo,
		View: viewInfo,
	}

	// 测试分析需要的字段
	var physicsTableFields []*types.FieldDefine
	required := sqlGen.AnalyzeRequiredFields([]UnsViewInfo{unsView}, physicsTableFields)

	// 验证需要的字段
	assert.Contains(t, required.FieldsToAdd, "double_1")
	assert.Contains(t, required.FieldsToAdd, "double_3")
	assert.Contains(t, required.FieldsToAdd, "long_1")

	// 验证字段类型映射
	assert.Equal(t, types.FieldType(types.FieldTypeDouble), required.FieldTypeMap["double_1"])
	assert.Equal(t, types.FieldType(types.FieldTypeDouble), required.FieldTypeMap["double_3"])
	assert.Equal(t, types.FieldType(types.FieldTypeLong), required.FieldTypeMap["long_1"])
}

func TestSQLGenerator_ExistingTable(t *testing.T) {
	sqlGen := NewSQLGenerator()

	// 模拟已存在的物理表字段
	physicsTableFields := []*types.FieldDefine{
		{Name: "tag", Type: types.FieldTypeLong},
		{Name: "timeStamp", Type: types.FieldTypeDatetime},
		{Name: "quality", Type: types.FieldTypeLong},
		{Name: "double_1", Type: types.FieldTypeDouble},
		{Name: "double_2", Type: types.FieldTypeDouble},
		{Name: "double_3", Type: types.FieldTypeDouble},
	}

	unsList := []UnsViewInfo{
		{
			Uns: &types.CreateTopicDto{
				Id:        1003,
				Alias:     "test_view_3",
				TableName: "uns_timeserial",
				Fields: []*types.FieldDefine{
					{Name: "tag", Type: types.FieldTypeLong, SystemField: base.OptionalTrue, Unique: base.OptionalTrue},
					{Name: "timeStamp", Type: types.FieldTypeDatetime, SystemField: base.OptionalTrue, Unique: base.OptionalTrue},
					{Name: "quality", Type: types.FieldTypeLong},
					{Name: "double1", Type: types.FieldTypeDouble},
					{Name: "double4", Type: types.FieldTypeDouble},
					{Name: "double5", Type: types.FieldTypeDouble},
				},
			},
			View: SimpleViewInfo{
				SrcTable: "uns_timeserial",
				Columns: []ViewColumnInfo{
					{ColumnName: "timeStamp"},
					{ColumnName: "quality"},
					{ColumnName: "double1", SourceColumn: "double_1", Expression: "1+double_1 as double1"},
					{ColumnName: "double2", SourceColumn: "double_2"},
					{ColumnName: "double3", SourceColumn: "double_3"},
				},
			},
		},
	}

	// 生成 SQL
	result := sqlGen.GenerateSyncSQLs(physicsTableFields, constants.SysFieldCreateTime, constants.QosField, unsList)

	// 验证结果
	//assert.False(t, result.HasErrors())
	//assert.Equal(t, 0, len(result.CreateTableSQL))  // 表已存在，不需要创建
	//assert.Greater(t, len(result.AlterTableSQL), 0) // 需要添加新字段
	//assert.Greater(t, len(result.CreateViewSQL), 0) // 需要创建视图
	//
	//t.Logf("Alter Table SQL:\n%s", result.AlterTableSQL[0])
	//t.Logf("Update Table SQL:\n%v", result.UpdateDataSQL)
	//t.Logf("Create View SQL:\n%s", result.CreateViewSQL[0])

	jsobs, _ := json.MarshalIndent(result, "", " ")
	t.Logf("genSql: %+v\n", string(jsobs))
}
func TestSQLGenerator_View_ChangeFieldName(t *testing.T) {
	sqlGen := NewSQLGenerator()

	// 模拟已存在的物理表字段
	physicsTableFields := []*types.FieldDefine{
		{Name: "tag", Type: types.FieldTypeLong},
		{Name: "timeStamp", Type: types.FieldTypeDatetime},
		{Name: "quality", Type: types.FieldTypeLong},
		{Name: "long_1", Type: types.FieldTypeLong},
		{Name: "double_1", Type: types.FieldTypeDouble},
		{Name: "double_2", Type: types.FieldTypeDouble},
		{Name: "double_3", Type: types.FieldTypeDouble},
	}

	unsList := []UnsViewInfo{
		{
			Uns: &types.CreateTopicDto{
				Id:        1003,
				Alias:     "test_view_3",
				TableName: "uns_timeserial",
				Fields: []*types.FieldDefine{
					{Name: "timeStamp", Type: types.FieldTypeDatetime, SystemField: base.OptionalTrue, Unique: base.OptionalTrue},
					{Name: "quality", Type: types.FieldTypeLong},
					{Name: "zz1", Type: types.FieldTypeInteger, Index: base.V2p("long_1")},
					{Name: "bb", Type: types.FieldTypeDouble, Index: base.V2p("double_1")},
				},
			},
			View: SimpleViewInfo{
				SrcTable: "uns_timeserial",
				Columns: []ViewColumnInfo{
					{ColumnName: "timeStamp"},
					{ColumnName: "quality"},
					{ColumnName: "zz", SourceColumn: "long_1"},
					{ColumnName: "bb", SourceColumn: "double_2"},
				},
			},
		},
	}

	// 生成 SQL
	result := sqlGen.GenerateSyncSQLs(physicsTableFields, constants.SysFieldCreateTime, constants.QosField, unsList)

	// 验证结果
	//assert.False(t, result.HasErrors())
	//assert.Equal(t, 0, len(result.CreateTableSQL))  // 表已存在，不需要创建
	//assert.Greater(t, len(result.AlterTableSQL), 0) // 需要添加新字段
	//assert.Greater(t, len(result.CreateViewSQL), 0) // 需要创建视图
	//
	//t.Logf("Alter Table SQL:\n%s", result.AlterTableSQL[0])
	//t.Logf("Update Table SQL:\n%v", result.UpdateDataSQL)
	//t.Logf("Create View SQL:\n%s", result.CreateViewSQL[0])

	jsobs, _ := json.MarshalIndent(result, "", " ")
	t.Logf("genSql: %+v\n", string(jsobs))
}
