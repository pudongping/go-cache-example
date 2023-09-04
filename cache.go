package cache

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
