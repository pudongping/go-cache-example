package cache

import (
	"log"
	"sync"
)

type Cache interface {
	// Set 设置/添加一个缓存，如果 key 存在，用新值覆盖旧值
	Set(key string, value interface{})
	// Get 通过 key 获取一个缓存值
	Get(key string) interface{}
	// Del 通过 key 删除一个缓存值
	Del(key string)
	// DelOldest 删除最“无用”的一个缓存值
	DelOldest()
	// Len 获取缓存已存在的记录数
	Len() int
}

// DefaultMaxBytes 默认允许占用的最大内存
const DefaultMaxBytes = 1 << 29

// safeCache 并发安全缓存
type safeCache struct {
	m     sync.RWMutex
	cache Cache

	nget int // 记录缓存获取次数
	nhit int // 记录缓存命中次数
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++
	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hit")
		sc.nhit++
	}

	return v
}

func (sc *safeCache) stat() *Stat {
	sc.m.RLock()
	defer sc.m.RUnlock()
	return &Stat{
		NHit: sc.nhit,
		NGet: sc.nget,
	}
}

// 方便查看统计数据
type Stat struct {
	NHit, NGet int
}
