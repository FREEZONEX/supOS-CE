package event

import "backend/internal/common/dto/grafana"

// BatchCreateDashboardEvent defines an event for batch creating Grafana dashboards.
type BatchCreateDashboardEvent struct {
	DashboardDtoList []*grafana.DashboardDto
}
