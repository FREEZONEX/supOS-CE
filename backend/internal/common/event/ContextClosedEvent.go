package event

import "backend/internal/svc"

type ContextClosedEvent struct {
	SvcContext *svc.ServiceContext
}
