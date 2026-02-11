package msg_consumer

import (
	"backend/internal/types"
	"strconv"
	"testing"
)
import "github.com/dgraph-io/ristretto/v2"

func TestFmtInt(t *testing.T) {
	ts := int64(2008743696706572307)
	t.Log("ts:", ts)
	t.Log("s16: " + strconv.FormatInt(ts, 16))
	t.Log("s36: " + strconv.FormatInt(ts, 36))

	t.Log(strconv.ParseInt("f9ew7mt5se8j", 36, 64))
}
func Test_ecache(t *testing.T) {

	cache, err := ristretto.NewCache(&ristretto.Config[string, *types.UnsDefinition]{
		NumCounters: 1e6,     // number of keys to track frequency of (1M).
		MaxCost:     1 << 28, // maximum cost of cache (256M).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	// set a value with a cost of 1
	u1 := &types.UnsDefinition{CreateTopicDto: types.CreateTopicDto{Id: 121, Alias: "a121", Name: "seqx"}}
	cache.Set("key", u1, 1)
	// wait for value to pass through buffers
	cache.Wait()

	cache.RemainingCost()

	// get value from cache
	value, found := cache.Get("key")
	if !found {
		panic("missing value")
	}
	t.Logf("%p.%p vs %p.%p, value = %v\n", u1, &u1.Lock, value, &value.Lock, value)

	// del value from cache
	cache.Del("key")
}
