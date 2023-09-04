package cache

type Value interface {
	// 返回占用的内存字节数
	Len() int
}
