package fifo

import (
	"container/list"

	"github.com/pudongping/go-cache-example"
)

// fifo 是一个 FIFO cache。它不是并发安全的。
type fifo struct {
	// 缓存最大的容量，单位字节；
	// groupcache 使用的是最大存放 entry 个数
	maxBytes int

	// 当一个 entry 从缓存中移除是调用该回调函数，默认为 nil
	// groupcache 中的 key 是任意的可比较类型；value 是 interface{}
	onEvicted func(key string, value interface{})

	// 已使用的字节数，只包括值，key 不算
	usedBytes int

	// list.List 是 Go 标准库提供的双向链表。通过这个数据结构存放具体的值，
	// 可以做到移动记录到队尾的时间复杂度是 O(1)，在队尾增加记录时间复杂度也是 O(1)，同时删除一条记录时间复杂度同样是 O(1)。
	ll *list.List

	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

// New 创建一个新的 Cache，如果 maxBytes 是 0，表示没有容量限制
func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// Set 往 Cache 尾部增加一个元素（如果已经存在，则放入尾部，并修改值）
func (f *fifo) Set(key string, value interface{}) {
	// 如果 key 存在，则更新对应节点的值，并将该节点移到队尾
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	// 往队尾添加新节点
	en := &entry{key, value}
	e := f.ll.PushBack(en)
	// 并添加 key 和节点的映射关系
	f.cache[key] = e

	f.usedBytes += en.Len()
	// 如果超过了设定的最大值 f.maxBytes 则移除“无用”的节点
	// 这里设定的是，如果 maxBytes = 0，表示不限内存，则不会进行移除操作
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

// Get 从 cache 中获取 key 对应的值，nil 表示 key 不存在
func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}

	return nil
}

// Del 从 cache 中删除 key 对应的记录
func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

// DelOldest 从 cache 中删除最旧的记录
func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

// Len 返回当前 cache 中的记录数
func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)

	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}
