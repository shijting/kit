package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type LRUCache[K comparable, V any] struct {
	mu   sync.Mutex
	list *list.List
	data map[K]*list.Element
	// 最大容量
	MaxCapacity int
}

func NewLRUCache[K comparable, V any](cap int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		list:        list.New(),
		data:        make(map[K]*list.Element),
		MaxCapacity: cap,
	}
}

type Item[K comparable, V any] struct {
	Key        K
	Value      V
	Expiration time.Time
}

func (l *LRUCache[K, V]) Get(key K) (value Item[K, V], ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var zeroVal Item[K, V]
	if ele, exists := l.data[key]; exists {
		// 懒删除，检查过期时间
		expiration := ele.Value.(Item[K, V]).Expiration
		if !expiration.IsZero() && time.Now().After(expiration) {
			// 过期则删除
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
}

// RemoveBack
func (l *LRUCache[K, V]) removeOldest() {
	l.mu.Lock()
	defer l.mu.Unlock()
	ele := l.list.Back()
	if ele != nil {
		l.removeElement(ele)
	}
}

func (l *LRUCache[K, V]) removeElement(ele *list.Element) {
	delete(l.data, ele.Value.(Item[K, V]).Key)
	l.list.Remove(ele)
}

func (l *LRUCache[K, V]) Print() {
	l.mu.Lock()
	defer l.mu.Unlock()

	ele := l.list.Front()
	for ele != nil {
		fmt.Println(ele.Value.(V))
		ele = ele.Next()
	}
}
