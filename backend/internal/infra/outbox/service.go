package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/repo"
)

type Service struct {
	store *repo.OutboxRepo
}

func New() *Service {
	return &Service{store: repo.NewOutboxRepo(context.Background())}
}

func (s *Service) Enqueue(ctx context.Context, eventType, aggregateType, aggregateID string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	eventID := fmt.Sprintf("%s:%s:%d", eventType, aggregateID, time.Now().UTC().UnixNano())
	return s.store.EnqueueOutbox(ctx, repo.OutboxEvent{
		EventID:       eventID,
		EventType:     eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       body,
	})
}

func NextRetryTime(attempts int, now time.Time) int64 {
	delay := 30 * time.Second
	switch attempts {
	case 0:
		delay = 0
	case 1:
		delay = 30 * time.Second
	case 2:
		delay = 2 * time.Minute
	case 3:
		delay = 10 * time.Minute
	default:
		delay = 30 * time.Minute
	}
	return now.Add(delay).UTC().UnixMilli()
}
