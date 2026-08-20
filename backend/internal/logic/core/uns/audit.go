package uns

import (
	"context"
	"strconv"
	"strings"

	auditdomain "backend/internal/domain/audit"
	"backend/internal/svc"
)

func recordUNSAudit(ctx context.Context, svcCtx *svc.ServiceContext, data map[string]any, businessType string) {
	if svcCtx == nil || svcCtx.App == nil || svcCtx.App.Audit == nil {
		return
	}
	id := unsAuditString(data, "id")
	name := firstUNSAuditString(data, "namespace", "name", "displayName")
	if id == "" && name == "" {
		return
	}
	detail := map[string]any{"name": name}
	if id != "" {
		detail["id"] = id
		detail["path"] = "/api/core/uns/" + id
	}
	svcCtx.App.Audit.Record(ctx, auditdomain.RecordInput{
		ScopeType:    auditdomain.ScopeTypePlatform,
		ResType:      auditdomain.ResTypeUNS,
		ResID:        id,
		ResName:      name,
		BusinessType: businessType,
		Detail:       detail,
	})
}

func firstUNSAuditString(data map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := unsAuditString(data, key); value != "" {
			return value
		}
	}
	return ""
}

func unsAuditString(data map[string]any, key string) string {
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
