# go-cache-example

进程内缓存示例。

## 常用的缓存淘汰算法

-  LRU（Least Recently Used）：最近最少使用。移除最长时间不被使用的对象。
-  LFU（Least Frequently Used）：最不经常使用。移除使用次数最少的对象。
-  FIFO（First In First Out）：先进先出。按对象进入缓存的顺序来移除它们。