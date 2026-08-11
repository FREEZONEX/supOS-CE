package events

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type Phase string

const (
	PhaseInTx        Phase = "inTx"
	PhaseAfterCommit Phase = "afterCommit"
	PhaseAsync       Phase = "async"
)

type Handler func(ctx context.Context, event any) error

type handlerEntry struct {
	order int
	name  string
	fn    Handler
}

type Bus struct {
	mu       sync.RWMutex
	handlers map[Phase]map[reflect.Type][]handlerEntry
}

var defaultBus = NewBus()

func NewBus() *Bus {
	return &Bus{handlers: map[Phase]map[reflect.Type][]handlerEntry{}}
}

func Default() *Bus {
	return defaultBus
}

func (b *Bus) OnInTx(sample any, order int, name string, fn Handler) {
	b.register(PhaseInTx, sample, order, name, fn)
}

func (b *Bus) OnAfterCommit(sample any, order int, name string, fn Handler) {
	b.register(PhaseAfterCommit, sample, order, name, fn)
}

func (b *Bus) OnAsync(sample any, order int, name string, fn Handler) {
	b.register(PhaseAsync, sample, order, name, fn)
}

func (b *Bus) Publish(ctx context.Context, phase Phase, event any) error {
	for _, entry := range b.entries(phase, event) {
		if err := entry.fn(ctx, event); err != nil {
			return fmt.Errorf("%s %s: %w", phase, entry.name, err)
		}
	}
	return nil
}

func (b *Bus) register(phase Phase, sample any, order int, name string, fn Handler) {
	if sample == nil || fn == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	typ := reflect.TypeOf(sample)
	if b.handlers[phase] == nil {
		b.handlers[phase] = map[reflect.Type][]handlerEntry{}
	}
	b.handlers[phase][typ] = append(b.handlers[phase][typ], handlerEntry{order: order, name: name, fn: fn})
	sort.SliceStable(b.handlers[phase][typ], func(i, j int) bool {
		return b.handlers[phase][typ][i].order < b.handlers[phase][typ][j].order
	})
}

func (b *Bus) entries(phase Phase, event any) []handlerEntry {
	if event == nil {
		return nil
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	entries := b.handlers[phase][reflect.TypeOf(event)]
	out := make([]handlerEntry, len(entries))
	copy(out, entries)
	return out
}
