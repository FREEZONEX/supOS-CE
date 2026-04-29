package FieldUtils

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"context"
	"testing"
)

func TestProcessFieldDefinesKeepsJsonbPayloadField(t *testing.T) {
	result, err := ProcessFieldDefines(
		context.Background(),
		types.SrcJdbcTypePostgresql,
		[]*types.FieldDefine{{Name: constants.JsonbField, Type: types.FieldTypeString}},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("ProcessFieldDefines returned error: %v", err)
	}
	if result == nil {
		t.Fatal("expected processed fields, got nil")
	}

	names := make([]string, 0, len(result.Fields))
	for _, field := range result.Fields {
		names = append(names, field.Name)
	}

	expected := []string{constants.SysFieldCreateTime, constants.JsonbField, constants.SysFieldID}
	if len(names) != len(expected) {
		t.Fatalf("expected fields %v, got %v", expected, names)
	}
	for i := range expected {
		if names[i] != expected[i] {
			t.Fatalf("expected fields %v, got %v", expected, names)
		}
	}
}
