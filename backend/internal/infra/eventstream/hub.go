package eventstream

import (
	"sync"
	"time"
)

type Event struct {
	Type        string         `json:"type"`
	Keys        []string       `json:"keys,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
	CreatedTime int64          `json:"createdTime"`
}

type Hub struct {
	mu          sync.Mutex
	nextID      int64
	subscribers map[int64]subscriber
}

type subscriber struct {
	ch    chan Event
	match func(Event) bool
}

var Default = NewHub()

func NewHub() *Hub {
	return &Hub{subscribers: map[int64]subscriber{}}
}

func Publish(event Event) {
	Default.Publish(event)
}

func Subscribe(match func(Event) bool) (<-chan Event, func()) {
	return Default.Subscribe(match)
}

func (h *Hub) Publish(event Event) {
	if h == nil {
		return
	}
	if event.CreatedTime == 0 {
		event.CreatedTime = time.Now().UTC().UnixMilli()
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, sub := range h.subscribers {
		if sub.match != nil && !sub.match(event) {
			continue
		}
		select {
		case sub.ch <- event:
		default:
			select {
			case <-sub.ch:
			default:
			}
			select {
			case sub.ch <- event:
			default:
			}
		}
	}
}

func (h *Hub) Subscribe(match func(Event) bool) (<-chan Event, func()) {
	if h == nil {
		ch := make(chan Event)
		close(ch)
		return ch, func() {}
	}
	h.mu.Lock()
	h.nextID++
	id := h.nextID
	ch := make(chan Event, 64)
	h.subscribers[id] = subscriber{ch: ch, match: match}
	h.mu.Unlock()

	cancel := func() {
		h.mu.Lock()
		sub, ok := h.subscribers[id]
		if ok {
			delete(h.subscribers, id)
			close(sub.ch)
		}
		h.mu.Unlock()
	}
	return ch, cancel
}

func ContainsKey(keys []string, key string) bool {
	for _, item := range keys {
		if item == key {
			return true
		}
	}
	return false
}
