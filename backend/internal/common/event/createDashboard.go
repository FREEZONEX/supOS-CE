package event

import (
	"context"
	"encoding/json"
	"log"
)

// CreateDashboardEvent defines an event for creating a Grafana dashboard record.
type CreateDashboardEvent struct {
	ApplicationEvent
	UUID        string
	UnsAlias    []string
	Name        string
	UserName    string
	Description string
}

func NewCreateDashboardEvent(ctx context.Context, UnsAlias []string, uuid, name, description, userName string) *CreateDashboardEvent {
	return &CreateDashboardEvent{
		ApplicationEvent: ApplicationEvent{Context: ctx},
		UnsAlias:         UnsAlias,
		UUID:             uuid,
		Name:             name,
		Description:      description,
		UserName:         userName,
	}
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
