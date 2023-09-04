package lru_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/pudongping/go-cache-example/lru"
)

// go test -run TestSet
func TestSet(t *testing.T) {
	is := is.New(t)

	cache := lru.New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1)

	cache.Del("k1")
	is.Equal(0, cache.Len())

	// cache.Set("k2", time.Now())
}

// go test -run TestOnEvicted
func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := lru.New(16, onEvicted)

	// 因为 k1 最近一直被访问，因此它不会被淘汰
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{}

	is.Equal(expected, keys)
	is.Equal(4, cache.Len())
}