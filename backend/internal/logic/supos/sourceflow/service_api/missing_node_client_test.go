package service_api

import (
	"testing"

	"backend/internal/common/constants"
)

func TestNormalizeMissingNodeFlowType(t *testing.T) {
	tests := []struct {
		name     string
		flowType string
		want     string
	}{
		{name: "empty defaults to source node-red", flowType: "", want: constants.FlowTypeNODERED},
		{name: "source constant", flowType: constants.FlowTypeNODERED, want: constants.FlowTypeNODERED},
		{name: "event constant", flowType: constants.FlowTypeEVENTFLOW, want: constants.FlowTypeEVENTFLOW},
		{name: "event compact alias", flowType: "eventflow", want: constants.FlowTypeEVENTFLOW},
		{name: "event numeric alias", flowType: "2", want: constants.FlowTypeEVENTFLOW},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeMissingNodeFlowType(tt.flowType); got != tt.want {
				t.Fatalf("normalizeMissingNodeFlowType(%q) = %q, want %q", tt.flowType, got, tt.want)
			}
		})
	}
}
