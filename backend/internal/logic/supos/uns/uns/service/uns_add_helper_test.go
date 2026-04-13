package service

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"context"
	"strings"
	"testing"
)

func TestInitParamsUnsRejectsDescriptionOverLimit(t *testing.T) {
	description := strings.Repeat("a", maxDescriptionLen+1)
	errTipMap := map[string]string{}
	dto := &types.CreateTopicDto{
		Index:       3,
		Name:        "machine_state",
		Alias:       "machine_state",
		PathType:    constants.PathTypeDir,
		Description: &description,
	}

	pathMap := initParamsUns(context.Background(), []*types.CreateTopicDto{dto}, errTipMap)

	if len(errTipMap) != 1 {
		t.Fatalf("expected one validation error, got %d: %#v", len(errTipMap), errTipMap)
	}
	if errTipMap[dto.GainBatchIndex()] == "" {
		t.Fatalf("expected error for batch index %s, got %#v", dto.GainBatchIndex(), errTipMap)
	}
	if len(pathMap[constants.PathTypeDir]) != 0 {
		t.Fatalf("invalid dto should not be accepted into path map: %#v", pathMap[constants.PathTypeDir])
	}
}
