package repo

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OutboxEvent struct {
	ID              int64           `gorm:"column:id;primaryKey" json:"id"`
	EventID         string          `gorm:"column:event_id;uniqueIndex:idx_sys_outbox_event_event_id" json:"eventId"`
	EventType       string          `gorm:"column:event_type" json:"eventType"`
	AggregateType   string          `gorm:"column:aggregate_type" json:"aggregateType"`
	AggregateID     string          `gorm:"column:aggregate_id" json:"aggregateId"`
	Payload         json.RawMessage `gorm:"column:payload;type:jsonb" json:"payload"`
	Status          string          `gorm:"column:status" json:"status"`
	Attempts        int             `gorm:"column:attempts" json:"attempts"`
	NextRetryTime   int64           `gorm:"column:next_retry_time" json:"nextRetryTime"`
	LastError       string          `gorm:"column:last_error" json:"lastError"`
	LockedBy        string          `gorm:"column:locked_by" json:"lockedBy"`
	LockedUntilTime int64           `gorm:"column:locked_until_time" json:"lockedUntilTime"`
	OnlyTime
}

func (OutboxEvent) TableName() string { return "sys_outbox_event" }

func (r *OutboxRepo) EnqueueOutbox(ctx context.Context, event OutboxEvent) error {
	now := time.Now().UTC().UnixMilli()
	if len(event.Payload) == 0 {
		event.Payload = json.RawMessage(`{}`)
	}
	event.Status = "pending"
	event.Attempts = 0
	event.NextRetryTime = now
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "event_id"}}, DoNothing: true}).
		Create(&event).Error
}

// ClaimOutbox atomically locks and returns a batch of pending events. The
// SELECT ... FOR UPDATE SKIP LOCKED self-join has no gorm builder equivalent,
// so it stays as raw SQL.
func (r *OutboxRepo) ClaimOutbox(ctx context.Context, worker string, limit int, lockMs int64) ([]OutboxEvent, error) {
	if limit <= 0 {
		limit = 10
	}
	now := time.Now().UTC().UnixMilli()
	var out []OutboxEvent
	err := r.db.WithContext(ctx).Raw(`
UPDATE sys_outbox_event
SET locked_by=?, locked_until_time=?, updated_time=?
WHERE id IN (
    SELECT id FROM sys_outbox_event
    WHERE status='pending'
      AND next_retry_time <= ?
      AND (locked_until_time=0 OR locked_until_time < ?)
    ORDER BY created_time, id
    LIMIT ?
    FOR UPDATE SKIP LOCKED
)
RETURNING id, event_id, event_type, aggregate_type, aggregate_id, payload, status, attempts, next_retry_time, last_error, locked_by, locked_until_time, created_time, updated_time`,
		worker, now+lockMs, repoTimeFromMilli(now), now, now, limit).Scan(&out).Error
	return out, err
}

func (r *OutboxRepo) ClaimOutboxByTypes(ctx context.Context, worker string, eventTypes []string, limit int, lockMs int64) ([]OutboxEvent, error) {
	if len(eventTypes) == 0 {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}
	now := time.Now().UTC().UnixMilli()
	var out []OutboxEvent
	err := r.db.WithContext(ctx).Raw(`
UPDATE sys_outbox_event
SET locked_by=?, locked_until_time=?, updated_time=?
WHERE id IN (
    SELECT id FROM sys_outbox_event
    WHERE status='pending'
      AND event_type IN ?
      AND next_retry_time <= ?
      AND (locked_until_time=0 OR locked_until_time < ?)
    ORDER BY created_time, id
    LIMIT ?
    FOR UPDATE SKIP LOCKED
)
RETURNING id, event_id, event_type, aggregate_type, aggregate_id, payload, status, attempts, next_retry_time, last_error, locked_by, locked_until_time, created_time, updated_time`,
		worker, now+lockMs, repoTimeFromMilli(now), eventTypes, now, now, limit).Scan(&out).Error
	return out, err
}

func (r *OutboxRepo) MarkOutboxDone(ctx context.Context, eventID string) error {
	now := time.Now().UTC().UnixMilli()
	return r.db.WithContext(ctx).Model(&OutboxEvent{}).
		Where("event_id = ?", eventID).
		Updates(touchValues(map[string]any{"status": "done", "locked_by": "", "locked_until_time": 0}, now)).Error
}

func (r *OutboxRepo) MarkOutboxRetry(ctx context.Context, eventID string, attempts int, nextRetryTime int64, lastError string) error {
	now := time.Now().UTC().UnixMilli()
	status := "pending"
	if attempts >= 5 {
		status = "dead"
	}
	return r.db.WithContext(ctx).Model(&OutboxEvent{}).
		Where("event_id = ?", eventID).
		Updates(touchValues(map[string]any{
			"status":            status,
			"attempts":          attempts,
			"next_retry_time":   nextRetryTime,
			"last_error":        lastError,
			"locked_by":         "",
			"locked_until_time": 0,
		}, now)).Error
}

type OutboxRepo struct{ db *gorm.DB }

func NewOutboxRepo(in any) *OutboxRepo { return &OutboxRepo{db: GetCommonConn(in)} }
