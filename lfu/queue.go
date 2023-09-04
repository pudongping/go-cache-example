package lfu

import (
	"container/heap"

	"github.com/pudongping/go-cache-example"
)

type entry struct {
	key    string
	value  interface{}
	weight int // 表示该 entry 在 queue 中权重（优先级），访问次数越多，权重越高
	index  int // 代表该 entry 在堆（heap）中的索引
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil // avoid memory leak
	en.index = -1  // for safety
	*q = old[0 : n-1]
	return en
}

// update modifies the weight and value of an entry in the queue.
func (q *queue) update(en *entry, value interface{}, weight int) {
	en.value = value
	en.weight = weight
	heap.Fix(q, en.index)
}
