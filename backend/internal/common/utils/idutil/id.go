package idutil

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake ID generator
const (
	defaultEpoch     = 1288834974657 // Default epoch (2010-11-04 01:42:54)
	workerIDBits     = 5
	datacenterIDBits = 5
	sequenceBits     = 12

	maxWorkerID     = -1 ^ (-1 << workerIDBits)
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	maxSequence     = -1 ^ (-1 << sequenceBits)

	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits

	// DefaultTimeOffset 默认时钟回拨容忍时间 (5分钟)
	DefaultTimeOffset = 300_000 // milliseconds
)

// Snowflake generates unique IDs using Snowflake algorithm
type Snowflake struct {
	mu               sync.Mutex
	lastTimestamp    int64
	workerID         int64
	datacenterID     int64
	sequence         int64
	epoch            int64
	isUseSystemClock bool
	timeOffset       int64
}

var (
	suposSnowflake *Snowflake
	snowflakeMu    sync.Mutex
)

// NextID generates next unique ID
func NextID() int64 {
	return GetSnowflake(nil, 0, 0, false, DefaultTimeOffset).NextID()
}

// GetSnowflake returns a singleton Snowflake instance
func GetSnowflake(epochDate *time.Time, workerId, dataCenterId int64, isUseSystemClock bool, timeOffset int64) *Snowflake {
	if suposSnowflake == nil {
		snowflakeMu.Lock()
		defer snowflakeMu.Unlock()
		if suposSnowflake == nil {
			epoch := int64(defaultEpoch)
			if epochDate != nil {
				epoch = epochDate.UnixMilli()
			}

			// Validate worker ID and datacenter ID
			if workerId < 0 || workerId > maxWorkerID {
				workerId = 0
			}
			if dataCenterId < 0 || dataCenterId > maxDatacenterID {
				dataCenterId = 0
			}

			suposSnowflake = &Snowflake{
				epoch:            epoch,
				workerID:         workerId,
				datacenterID:     dataCenterId,
				sequence:         0,
				isUseSystemClock: isUseSystemClock,
				timeOffset:       timeOffset,
			}
		}
	}
	return suposSnowflake
}

// NextId generates a new unique ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := s.getCurrentTimestamp()

	// Handle clock moving backwards
	if timestamp < s.lastTimestamp {
		diff := s.lastTimestamp - timestamp
		if diff > s.timeOffset {
			panic(fmt.Sprintf("Clock moved backwards. Refusing to generate id for %d milliseconds", diff))
		}
		// Wait until clock catches up
		time.Sleep(time.Duration(diff) * time.Millisecond)
		timestamp = s.getCurrentTimestamp()
	}

	if timestamp == s.lastTimestamp {
		// Same millisecond, increment sequence
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// Sequence overflow, wait for next millisecond
			for timestamp <= s.lastTimestamp {
				timestamp = s.getCurrentTimestamp()
			}
		}
	} else {
		// New millisecond, reset sequence
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// Generate ID
	id := ((timestamp - s.epoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// getCurrentTimestamp returns current timestamp in milliseconds
func (s *Snowflake) getCurrentTimestamp() int64 {
	if s.isUseSystemClock {
		// Use System.currentTimeMillis() equivalent for better performance
		return time.Now().UnixMilli()
	}
	return time.Now().UnixMilli()
}
