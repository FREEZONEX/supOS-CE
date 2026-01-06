package diskqueue

import (
	"os"
	"testing"
	"time"
)

func TestDiskQueue_WriteAndRead(t *testing.T) {
	dq := newQueue(t)
	//defer dq.Close()
	data := make([]byte, 1024)

	for i := 0; i < 30; i++ {
		for j := 0; j < len(data); j++ {
			data[j] = byte(i)
		}
		err := dq.Put(data)
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(1 * time.Second)
	go func() {
		tk := time.NewTicker(2 * time.Second)
		run := true
		for i := 0; run; i++ {
			select {
			case <-tk.C:
				tk.Reset(2 * time.Second)
			case msg := <-dq.ReadChan():
				t.Logf("msg: %v - %v\n", len(msg), msg[0])
				if i == 5 {
					run = false
					break
				}
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
func TestDiskQueue_ReadOfflineQueue(t *testing.T) {
	dq := newQueue(t)
	defer dq.Close()
	go func() {
		tk := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-tk.C:
				tk.Reset(2 * time.Second)
			case msg := <-dq.ReadChan():
				t.Logf("msg: %v - %v\n", len(msg), msg[0])
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
func newQueue(t *testing.T) Interface {
	l := NewTestLogger(t)
	dqName := "bench_disk_queue_put"
	dir := "F:/var/queue"
	err := os.MkdirAll(dir, 666)
	if err != nil {
		panic(err)
	}
	return New(dqName, dir, 20*(1<<10), 0, 1<<10, 2500, 2*time.Second, l)
}
