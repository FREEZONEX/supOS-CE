package service_api

import (
	"strings"

	"backend/internal/common/constants"
	"backend/internal/svc"
	noderedclient "backend/share/clients/nodered"
)

func missingNodeRedClient(svcCtx *svc.ServiceContext, flowType string) *noderedclient.Client {
	if svcCtx == nil {
		return nil
	}
	if normalizeMissingNodeFlowType(flowType) == constants.FlowTypeEVENTFLOW && svcCtx.EventNodeRed != nil {
		return svcCtx.EventNodeRed
	}
	return svcCtx.SourceNodeRed
}

func normalizeMissingNodeFlowType(flowType string) string {
	switch strings.ToLower(strings.TrimSpace(flowType)) {
	case constants.FlowTypeEVENTFLOW, "eventflow", "event_flow", "2":
		return constants.FlowTypeEVENTFLOW
	default:
		return constants.FlowTypeNODERED
	}
}
