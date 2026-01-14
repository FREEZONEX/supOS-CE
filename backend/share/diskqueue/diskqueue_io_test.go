package diskqueue

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
				//if i == 5 {
				//	run = false
				//	break
				//}
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
func TestDiskQueue_ReadOfflineQueue(t *testing.T) {
	dq := newQueue(t)
	defer dq.Close()
	go func() {
		tk := time.NewTicker(time.Second)
		for {
			select {
			case <-tk.C:
			case msg := <-dq.ReadChan():
				t.Logf("msg: %v - %v\n", len(msg), msg[0])
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
func newQueue(t *testing.T) Interface {
	_, q := newQueueWith(t, 10*1024, 1, 1<<10)
	return q
}
func newQueueWith(t *testing.T, maxBytesPerFile int64, minMsgSize, maxMsgSize int32) (string, Interface) {
	dqName := "unstest"
	dir := "F:/var/queue"
	err := os.MkdirAll(dir, 666)
	if err != nil {
		panic(err)
	}
	return dir, New(dqName, dir, maxBytesPerFile, minMsgSize, maxMsgSize, 2500, 2*time.Second, t.Logf, t.Errorf)
}

func TestLargeMessageBoundary(t *testing.T) {
	// Use smaller sizes to test the same behavior more efficiently
	// 5KB file limit, 4KB max message (same 10:8 ratio as 50MB:40MB in production)
	maxBytesPerFile := int64(5 * 1024)
	maxMsgSize := int32(4 * 1024)

	dir, dq := newQueueWith(t, maxBytesPerFile, 1, maxMsgSize)
	defer dq.Close()

	// Create messages that will cause multiple rotations
	largeMsg := make([]byte, 4000) // ~4KB message
	for i := 0; i < 15; i++ {      // ~60KB total, should rotate cleanly across multiple files
		err := dq.Put(largeMsg)
		Nil(t, err)
	}

	// Read all messages back
	for i := 0; i < 15; i++ {
		msg := <-dq.ReadChan()
		Equal(t, len(largeMsg), len(msg))
	}

	// Verify no .bad files were created
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".bad") {
			t.Fatalf("Found corrupted file: %s", path)
		}
		return err
	})
}

// TestReadCurrentWriteFile verifies that when reading the current write file,
// the reader doesn't try to rotate past the write file when reaching maxBytesPerFileRead
func TestReadCurrentWriteFile(t *testing.T) {

	// Small file limit to trigger boundary easily
	maxBytesPerFile := int64(1024)
	dir, dq := newQueueWith(t, maxBytesPerFile, 4, 1024)
	defer dq.Close()

	// Write messages up to the file limit
	msg := []byte("test message")
	for i := 0; i < 60; i++ { // Enough to fill first file and start second
		err := dq.Put(msg)
		Nil(t, err)
	}

	// Read all messages back - this tests reading from current write file
	// without trying to advance past it
	for i := 0; i < 60; i++ {
		readMsg := <-dq.ReadChan()
		Equal(t, msg, readMsg)
	}

	// Verify no .bad files were created
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".bad") {
			t.Fatalf("Found corrupted file: %s", path)
		}
		return err
	})
}
