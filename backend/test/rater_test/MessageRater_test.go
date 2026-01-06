package main

import (
	"fmt"
	"testing"
	"time"
)

func TestMsgRater(t *testing.T) {
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
}

func TestMsgCounter(t *testing.T) {
	counter := &MessageCounter{}

	// 启动统计协程
	go counter.Start()

	// 模拟消息接收
	for i := 0; i < 10000; i++ {
		go func() {
			for {
				counter.Increment()
				time.Sleep(time.Millisecond * 10)
			}
		}()
	}
	go func() {
		ticker := time.NewTicker(time.Millisecond * 1000)
		for range ticker.C {
			fmt.Printf("QPS: %d\n", counter.GetQPS())
		}
	}()
	time.Sleep(time.Second * 10)
}
