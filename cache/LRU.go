package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// LRUCache 是一个带有固定最大容量的最近最少使用（LRU）缓存实现。
// 它支持键值对，其中键的类型为 K（可比较的），值的类型为 V。
type LRUCache[K comparable, V any] struct {
	mu   sync.Mutex
	list *list.List
	data map[K]*list.Element
	// Maximum capacity of the LRUCache
	MaxCapacity int
}

// NewLRUCache 创建一个具有指定最大容量的新 LRUCache 实例。
func NewLRUCache[K comparable, V any](cap int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		list:        list.New(),
		data:        make(map[K]*list.Element),
		MaxCapacity: cap,
	}
}

// Item represents a key-value pair with an optional expiration time.
type Item[K comparable, V any] struct {
	Key        K
	Value      V
	Expiration time.Time // 过期时间，零表示永不过期
}

// Get 从 LRUCache 中检索与指定键关联的值。
// 如果找到键，则返回值和 true，否则返回零值和 false。
func (l *LRUCache[K, V]) Get(key K) (value Item[K, V], ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var zeroVal Item[K, V]
	if ele, exists := l.data[key]; exists {
		// Lazy deletion
		expiration := ele.Value.(Item[K, V]).Expiration
		if !expiration.IsZero() && time.Now().After(expiration) {
			l.removeElement(ele)
			return zeroVal, false
		}

		l.list.MoveToFront(ele)
		return ele.Value.(Item[K, V]), true
	}
	if len(l.data) > l.MaxCapacity {
		l.removeOldest()
	}
	return zeroVal, false
}

// Set 在 LRUCache 中添加或更新与指定键关联的值。
// 还允许为键值对设置可选的过期时间。 如果过期时间为0 表示永不过期。
func (l *LRUCache[K, V]) Set(key K, value V, expiration time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	newEle := Item[K, V]{Key: key, Value: value}

	if expiration > 0 {
		newEle.Expiration = time.Now().Add(expiration)
	}

	if v, ok := l.data[key]; ok {
		v.Value = newEle
		l.list.MoveToFront(v)
	} else {
		l.data[key] = l.list.PushFront(newEle)
	}
	if len(l.data) > l.MaxCapacity {
		l.removeOldest()
	}
}

// RemoveBack
func (l *LRUCache[K, V]) removeOldest() {
	ele := l.list.Back()
	if ele != nil {
		l.removeElement(ele)
	}
}

func (l *LRUCache[K, V]) removeElement(ele *list.Element) {
	delete(l.data, ele.Value.(Item[K, V]).Key)
	l.list.Remove(ele)
}

// Print 打印 LRUCache 元素值，按最近访问的顺序排列。
func (l *LRUCache[K, V]) Print() {
	l.mu.Lock()
	defer l.mu.Unlock()

	ele := l.list.Front()
	for ele != nil {
		fmt.Println(ele.Value.(V))
		ele = ele.Next()
	}
}
