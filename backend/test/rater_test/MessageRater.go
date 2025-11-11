package main

import (
	"sync/atomic"
	"time"
)

/*
*

	   counter := newMessageRater(4)

		// 启动统计协程
		go counter.Start()

		// 模拟消息接收
		for i := 0; i < 10000; i++ {
			go func() {
				for {
					counter.Add(i % 2)
					time.Sleep(time.Millisecond * 10)
				}
			}()
		}
		go func() {
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				fmt.Printf("QPS-0: %d, QPS-1: %d\n", counter.GetQPS(0), counter.GetQPS(1))
			}
		}()

		time.Sleep(time.Second * 10)

		counter.Shutdown()
		fmt.Println("已停止!")

		time.Sleep(time.Second * 5)
*/
type messageRater struct {
	ticker *time.Ticker
	data   []uint64
}

func newMessageRater(size int) messageRater {
	return messageRater{data: make([]uint64, size)}
}
func (m *messageRater) Add(i int) {
	if m.ticker != nil {
		atomic.AddUint64(&m.data[i<<1], 1)
	}
}
func (m *messageRater) Shutdown() {
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = nil
	}
}

func (m *messageRater) Start() {
	m.ticker = time.NewTicker(time.Second)
	N := len(m.data)
	for range m.ticker.C {
		if m.ticker == nil {
			break
		}
		for i := 0; i < N; i += 2 {
			m.data[i+1] = atomic.SwapUint64(&m.data[i], 0)
		}
	}
}
func (m *messageRater) GetQPS(i int) (rs uint64) {
	rs = atomic.LoadUint64(&m.data[(i<<1)+1])
	return rs
}
