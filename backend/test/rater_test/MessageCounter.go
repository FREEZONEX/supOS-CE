package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type MessageCounter struct {
	prev  uint64
	count uint64
	total uint64
}

func (m *MessageCounter) Increment() {
	atomic.AddUint64(&m.count, 1)
	atomic.AddUint64(&m.total, 1)
}

func (m *MessageCounter) Start() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("~~QPS: %d\n", m.count)
		m.prev = m.count
		m.count = 0
		//atomic.SwapUint64(&m.count, 0)
	}
}
func (m *MessageCounter) GetQPS() uint64 {
	return m.prev
}
