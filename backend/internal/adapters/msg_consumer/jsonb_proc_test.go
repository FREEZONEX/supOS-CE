package msg_consumer

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"context"
	"testing"
)

func newJsonbUnsDefinitionForTest() *types.UnsDefinition {
	jsonbType := constants.JsonbType
	unique := true
	return &types.UnsDefinition{CreateTopicDto: types.CreateTopicDto{
		Alias:     "jsonb_topic",
		DataType:  &jsonbType,
		DataSrcID: types.SrcJdbcTypePostgresql.Id(),
		Fields: []*types.FieldDefine{
			{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime},
			{Name: constants.JsonbField, Type: types.FieldTypeString},
			{Name: constants.SysFieldID, Type: types.FieldTypeLong, Unique: &unique},
		},
	}}
}

func TestProcDataJsonbStoresRawPayloadInJsonbField(t *testing.T) {
	rawPayload := `{"temperature":25,"status":"ok"}`
	data, err := parseJsonList([]byte(rawPayload))
	if err != nil {
		t.Fatalf("parseJsonList returned error: %v", err)
	}

	list, errMsg := procData(context.WithValue(context.Background(), "payload", rawPayload), newJsonbUnsDefinitionForTest(), data)
	if errMsg != "" {
		t.Fatalf("procData returned error message: %s", errMsg)
	}
	if len(list) != 1 {
		t.Fatalf("expected one jsonb row, got %d: %#v", len(list), list)
	}
	if _, ok := list[0][""]; ok {
		t.Fatalf("jsonb payload should not be stored under an empty field: %#v", list[0])
	}
	if got := list[0][constants.JsonbField]; got != rawPayload {
		t.Fatalf("expected raw payload in %s, got %s", constants.JsonbField, got)
	}
	if list[0][constants.SysFieldCreateTime] == "" {
		t.Fatalf("expected timestamp field to be populated: %#v", list[0])
	}
}

func TestProcDataJsonbFallbackDoesNotWrapSingleObjectInArray(t *testing.T) {
	data := []map[string]string{{"temperature": "25", "status": `"ok"`}}

	list, errMsg := procData(context.Background(), newJsonbUnsDefinitionForTest(), data)
	if errMsg != "" {
		t.Fatalf("procData returned error message: %s", errMsg)
	}
	if len(list) != 1 {
		t.Fatalf("expected one jsonb row, got %d: %#v", len(list), list)
	}
	if got := list[0][constants.JsonbField]; got == "" || got[0] == '[' {
		t.Fatalf("expected single object json payload, got %s", got)
	}
}
