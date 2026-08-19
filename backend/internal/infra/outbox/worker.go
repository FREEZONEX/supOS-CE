package outbox

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"backend/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type EventHandler func(ctx context.Context, event repo.OutboxEvent) error

type Worker struct {
	store       *repo.OutboxRepo
	handlers    map[string]EventHandler
	workerID    string
	batchSize   int
	lockTimeout time.Duration
	interval    time.Duration

	mu     sync.Mutex
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorker(ctx context.Context) *Worker {
	host, _ := os.Hostname()
	host = strings.TrimSpace(host)
	if host == "" {
		host = "backend"
	}
	return &Worker{
		store:       repo.NewOutboxRepo(ctx),
		handlers:    map[string]EventHandler{},
		workerID:    fmt.Sprintf("%s-%d", host, os.Getpid()),
		batchSize:   10,
		lockTimeout: 2 * time.Minute,
		interval:    2 * time.Second,
	}
}

func (w *Worker) Register(eventType string, handler EventHandler) {
	eventType = strings.TrimSpace(eventType)
	if w == nil || eventType == "" || handler == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[eventType] = handler
}

func (w *Worker) Start(ctx context.Context) {
	if w == nil {
		return
	}
	w.mu.Lock()
	if w.cancel != nil {
		w.mu.Unlock()
		return
	}
	runCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.wg.Add(1)
	w.mu.Unlock()

	go w.loop(runCtx)
}

func (w *Worker) Stop() {
	if w == nil {
		return
	}
	w.mu.Lock()
	cancel := w.cancel
	w.cancel = nil
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	w.wg.Wait()
}

func (w *Worker) loop(ctx context.Context) {
	defer w.wg.Done()
	for {
		claimed, err := w.drain(ctx)
		if err != nil {
			logx.WithContext(ctx).Errorf("outbox worker drain failed: %v", err)
		}
		if claimed > 0 {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(w.interval):
		}
	}
}

func (w *Worker) drain(ctx context.Context) (int, error) {
	eventTypes, handlers := w.snapshotHandlers()
	if len(eventTypes) == 0 {
		return 0, nil
	}
	events, err := w.store.ClaimOutboxByTypes(ctx, w.workerID, eventTypes, w.batchSize, int64(w.lockTimeout/time.Millisecond))
	if err != nil {
		return 0, err
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].CreatedTime.Equal(events[j].CreatedTime) {
			return events[i].ID < events[j].ID
		}
		return events[i].CreatedTime.Before(events[j].CreatedTime)
	})
	for _, event := range events {
		handler := handlers[event.EventType]
		if handler == nil {
			continue
		}
		w.handle(ctx, event, handler)
	}
	return len(events), nil
}

func (w *Worker) snapshotHandlers() ([]string, map[string]EventHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	handlers := make(map[string]EventHandler, len(w.handlers))
	eventTypes := make([]string, 0, len(w.handlers))
	for eventType, handler := range w.handlers {
		handlers[eventType] = handler
		eventTypes = append(eventTypes, eventType)
	}
	sort.Strings(eventTypes)
	return eventTypes, handlers
}

func (w *Worker) handle(ctx context.Context, event repo.OutboxEvent, handler EventHandler) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err := fmt.Errorf("panic: %v\n%s", recovered, debug.Stack())
			w.markRetry(ctx, event, err)
		}
	}()
	if err := handler(ctx, event); err != nil {
		w.markRetry(ctx, event, err)
		return
	}
	if err := w.store.MarkOutboxDone(ctx, event.EventID); err != nil {
		logx.WithContext(ctx).Errorf("outbox mark done failed eventID=%s err=%v", event.EventID, err)
	}
}

func (w *Worker) markRetry(ctx context.Context, event repo.OutboxEvent, err error) {
	attempts := event.Attempts + 1
	next := NextRetryTime(attempts, time.Now().UTC())
	message := ""
	if err != nil {
		message = err.Error()
	}
	if markErr := w.store.MarkOutboxRetry(ctx, event.EventID, attempts, next, message); markErr != nil {
		logx.WithContext(ctx).Errorf("outbox mark retry failed eventID=%s err=%v", event.EventID, markErr)
	}
	logx.WithContext(ctx).Errorf("outbox handler failed eventID=%s eventType=%s attempts=%d err=%v", event.EventID, event.EventType, attempts, err)
}
