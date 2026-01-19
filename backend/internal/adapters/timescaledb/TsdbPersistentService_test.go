package timescaledb

import (
	"backend/internal/types"
	"testing"
)

func TestSave(t *testing.T) {
	db := newTsdbPersistentService("postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if db == nil {
		t.Fatal("NoDB")
	}
	unsInfo := &types.CreateTopicDto{
		Id:        17,
		Alias:     "test_ds",
		TableName: "uns_timeserial",
		Fields: procFields([]*types.FieldDefine{
			createField("value", types.FieldTypeFloat, nil),
			createField("desc", types.FieldTypeString, nil),
		}),
	}
	t.Log("unsInfo", unsInfo)

	err := db.Save([]types.UnsInfo{unsInfo})
	if err != nil {
		t.Fatal(err)
	}
}
