package msg_consumer

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
	"strconv"
	"testing"
	"time"
)

func Test_filterMsgByUns(t *testing.T) {
	unsInfo := types.CreateTopicDto{
		Id:        1,
		Alias:     "test_uns",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			{
				Name: "temperature",
				Type: types.FieldTypeFloat,
			},
			{
				Name:   "state",
				Type:   types.FieldTypeString,
				MaxLen: base.V2p(5),
			},
		}),
		Timestamps: [2]int64{0, 0},
	}
	dataList := []map[string]string{
		{constants.SysFieldCreateTime: strconv.FormatInt(time.Now().UnixMilli(), 10), "temperature": "20", "state": "1"},
		{constants.SysFieldCreateTime: strconv.FormatInt(time.Now().UnixMilli(), 10), "temperature": "21", "state": "123456"},
	}
	errMsg := filterMsgByUns(&types.UnsDefinition{CreateTopicDto: unsInfo}, dataList)
	t.Log(errMsg)
	t.Logf("%+v", dataList)
}
func createField(name, fieldType string) *types.FieldDefine {
	return &types.FieldDefine{
		Name: name,
		Type: fieldType,
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
