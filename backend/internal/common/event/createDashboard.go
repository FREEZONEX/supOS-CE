package event

import (
	"encoding/json"
	"log"
)

// CreateDashboardEvent defines an event for creating a Grafana dashboard record.
type CreateDashboardEvent struct {
	UUID        string
	Name        string
	Description string
}

// String implements the fmt.Stringer interface to provide a JSON representation.
func (e *CreateDashboardEvent) String() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Printf("failed to marshal CreateDashboardEvent to json: %v", err)
		return ""
	}
	return string(bytes)
}
