package service

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"testing"
)

func TestSetJsonFieldsFallbackUsesFieldsWhenProtocolJsonFieldsMissing(t *testing.T) {
	fields := []*types.FieldDefine{{Name: constants.JsonbField, Type: types.FieldTypeString}}
	detail := &types.InstanceDetail{Fields: fields}

	setJsonFieldsFallback(detail)

	if len(detail.JsonFields) != 1 || detail.JsonFields[0].Name != constants.JsonbField {
		t.Fatalf("expected jsonFields to fall back to fields, got %#v", detail.JsonFields)
	}
}

func TestSetJsonFieldsFallbackKeepsExistingJsonFields(t *testing.T) {
	existing := []*types.FieldDefine{{Name: "payload", Type: types.FieldTypeString}}
	detail := &types.InstanceDetail{
		Fields:     []*types.FieldDefine{{Name: constants.JsonbField, Type: types.FieldTypeString}},
		JsonFields: existing,
	}

	setJsonFieldsFallback(detail)

	if len(detail.JsonFields) != 1 || detail.JsonFields[0].Name != "payload" {
		t.Fatalf("expected existing jsonFields to be preserved, got %#v", detail.JsonFields)
	}
}
