package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/goburrow/cache"
	"github.com/karlseguin/ccache/v2"
)

func TestCCache(t *testing.T) {
	// 创建一个最大容量为 1000 的缓存
	cache := ccache.New(ccache.Configure().MaxSize(1000))
	defer cache.Stop()

	key := "my-key"
	duration := 5 * time.Second
	// 设置初始值，过期时间为 5 秒
	cache.Set(key, "value-1", duration)

	// 模拟访问
	time.Sleep(2 * time.Second)

	t.Log(cache.Get(key + "_none"))
	cache.Set(key+"_none", nil, duration)
	t.Log(cache.Get(key + "_none"))

	item := cache.Get(key)
	if item != nil && !item.Expired() {
		fmt.Println("访问成功:", item.Value(), ", ttl(s): ", item.TTL().Seconds())

		// 【关键步骤】模拟 expireAfterAccess：
		// 每次成功访问后，将过期时间向后顺延 5 秒
		item.Extend(duration)
		fmt.Println("已刷新过期时间")
	}

	// 再睡 4 秒 (总共 6 秒，如果不刷新早就过期了，但刷新了所以还在)
	time.Sleep(4 * time.Second)

	item = cache.Get(key)
	if item != nil && !item.Expired() {
		fmt.Println("依然有效:", item.Value())
	} else if item != nil {
		fmt.Println("已过期:", item.Value())
	} else {
		fmt.Println("已过期")
	}
}
func TestGuavaCache(t *testing.T) {
	load := func(k cache.Key) (cache.Value, error) {
		fmt.Printf("loading %v\n", k)
		time.Sleep(500 * time.Millisecond)
		return fmt.Sprintf("%d-%d", k, time.Now().Unix()), nil
	}
	remove := func(k cache.Key, v cache.Value) {
		fmt.Printf("removed %v (%v)\n", k, v)
	}
	// Create a new cache
	c := cache.NewLoadingCache(load,
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Second),
		cache.WithRefreshAfterWrite(20*time.Second),
		cache.WithRemovalListener(remove),
	)
	getTicker := time.Tick(2 * time.Second)
	reportTicker := time.Tick(30 * time.Second)
	for {
		select {
		case <-getTicker:
			k := 1 //rand.Intn(100)
			v, _ := c.Get(k)
			fmt.Printf("get %v: %v\n", k, v)
		case <-reportTicker:
			st := cache.Stats{}
			c.Stats(&st)
			fmt.Printf("total: %d, hits: %d (%.2f%%), misses: %d (%.2f%%), evictions: %d, load: %s (%s)\n",
				st.RequestCount(), st.HitCount, st.HitRate()*100.0, st.MissCount, st.MissRate()*100.0,
				st.EvictionCount, st.TotalLoadTime, st.AverageLoadPenalty())
		}
	}
}
