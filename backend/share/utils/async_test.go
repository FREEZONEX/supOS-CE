package utils

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestGo_ExecutesFunction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	var called bool
	Go(context.Background(), func() {
		called = true
		wg.Done()
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("utils.Go did not execute function")
	}

	if !called {
		t.Fatal("function was not called")
	}
}

func TestGo_RecoversPanic(t *testing.T) {
	// 只要不 panic 把整个测试进程搞崩，就说明 recover 生效
	done := make(chan struct{})
	Go(context.Background(), func() {
		defer close(done)
		panic("intentional panic")
	})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("utils.Go did not finish after panic")
	}
}
