package idgen

import (
	"hash/fnv"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Cloud-compatible 53-bit snowflake ID: 42-bit millisecond timestamp,
// 7-bit machine ID and 4-bit sequence.
const epochMilli = int64(1609430400000)

type Snowflake struct {
	machineID int64
	sequence  int64
	lastMilli int64
	mu        sync.Mutex
}

func NewSnowflake(machineID int64) *Snowflake {
	return &Snowflake{
		machineID: (machineID & 0x7F) << 4,
		lastMilli: nowMilli(),
	}
}

func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := nowMilli()
	if current < s.lastMilli {
		current = s.lastMilli
	}
	if current == s.lastMilli {
		s.sequence++
		if s.sequence > 15 {
			for {
				runtime.Gosched()
				current = nowMilli()
				if current > s.lastMilli {
					break
				}
			}
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}
	s.lastMilli = current
	return (((current - epochMilli) & 0x3FFFFFFFFFF) << 11) | s.machineID | s.sequence
}

var defaultSnowflake = NewSnowflake(defaultMachineID())

func NextID() int64 {
	return defaultSnowflake.NextID()
}

func nowMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func defaultMachineID() int64 {
	if raw := os.Getenv("SNOWFLAKE_NODE_ID"); raw != "" {
		if value, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return value
		}
	}
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "edge-backend"
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(hostname))
	return int64(h.Sum64() & 0x7F)
}
