package fifo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/pudongping/go-cache-example/fifo"
)

// go test -run TestSetGet
func TestSetGet(t *testing.T) {
	is := is.New(t)

	cache := fifo.New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1)

	cache.Del("k1")
	is.Equal(0, cache.Len()) // expect to be the same

	// cache.Set("k2", time.Now())
}

// go test -run TestOnEvicted
func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := fifo.New(16, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{}

	is.Equal(expected, keys) // 不相等
	is.Equal(4, cache.Len())
}
