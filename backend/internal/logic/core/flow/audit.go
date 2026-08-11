package flow

import (
	"context"
	"strconv"
	"strings"

	auditdomain "backend/internal/domain/audit"
	"backend/internal/svc"
)

func recordFlowAudit(ctx context.Context, svcCtx *svc.ServiceContext, data map[string]any, businessType string) {
	if svcCtx == nil || svcCtx.App == nil || svcCtx.App.Audit == nil {
		return
	}
	id := flowAuditString(data, "id")
	name := firstFlowAuditString(data, "name", "flowName")
	if id == "" && name == "" {
		return
	}
	flowType := firstFlowAuditString(data, "flowType")
	nodeType := firstFlowAuditString(data, "nodeType")
	resType := auditdomain.ResTypeSourceFlow
	if strings.EqualFold(flowType, "event") || flowType == "2" {
		resType = auditdomain.ResTypeEventFlow
	}
	detail := map[string]any{
		"name":     name,
		"flowType": flowType,
		"nodeType": nodeType,
	}
	if id != "" {
		detail["id"] = id
		detail["path"] = "/api/core/flows/" + id
	}
	svcCtx.App.Audit.Record(ctx, auditdomain.RecordInput{
		ScopeType:    auditdomain.ScopeTypePlatform,
		ResType:      resType,
		ResID:        id,
		ResName:      name,
		BusinessType: businessType,
		Detail:       detail,
	})
}

func firstFlowAuditString(data map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := flowAuditString(data, key); value != "" {
			return value
		}
	}
	return ""
}

func flowAuditString(data map[string]any, key string) string {
	value, ok := data[key]
	if !ok {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case int:
		return strconv.Itoa(typed)
	case int32:
		return strconv.FormatInt(int64(typed), 10)
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return ""
	}
}
