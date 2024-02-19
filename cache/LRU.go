package cache

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"github.com/shijting/kit/option"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrCacheNotFound = errors.New("cache not found")
)

// LRUCache 是一个基于本地内存带有固定最大容量的最近最少使用（LRU）缓存实现。
// 它支持键值对，其中键的类型为 K（comparable），值的类型为 any。
type LRUCache[K comparable, V any] struct {
	mu   sync.RWMutex
	list *list.List
	data map[K]*list.Element
	// Maximum capacity of the LRUCache
	MaxCapacity          int
	gcRandomDeletionStep int
	// 用于删除过期元素的定时器触发间隔
	gcInterval time.Duration
}

// NewLRUCache 创建一个具有指定最大容量的新 LRUCache 实例。
// 如果未指定任何选项，则使用默认选项。
func NewLRUCache[K comparable, V any](cap int, opts ...option.Option[LRUCache[K, V]]) *LRUCache[K, V] {
	cache := &LRUCache[K, V]{
		list:                 list.New(),
		data:                 make(map[K]*list.Element),
		MaxCapacity:          cap,
		gcRandomDeletionStep: 100,
		gcInterval:           time.Second,
	}
	option.Options[LRUCache[K, V]](opts).Apply(cache)
	go cache.startGC()
	return cache
}

// WithGCRandomDeletionStep 设置随机删除步长。
// step: 每次触发垃圾回收时，最多删除的元素数量, 该值不能小于1。
func WithGCRandomDeletionStep(step int) option.Option[LRUCache[string, any]] {
	return func(t *LRUCache[string, any]) {
		t.gcRandomDeletionStep = step
	}
}

// WithGCInterval 设置垃圾回收间隔。
// 如果interval设置为0，则禁用垃圾回收。
func WithGCInterval[K comparable, V any](interval time.Duration) option.Option[LRUCache[K, V]] {
	return func(t *LRUCache[K, V]) {
		t.gcInterval = interval
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
func (l *LRUCache[K, V]) Get(ctx context.Context, key K) (value Item[K, V], ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
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
func (l *LRUCache[K, V]) Set(ctx context.Context, key K, value V, expiration time.Duration) error {
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
	return nil
}

// Delete 从 LRUCache 中删除与指定键关联的值。
// 如果键不存在，则返回ErrCacheNotFound错误。
func (l *LRUCache[K, V]) Delete(ctx context.Context, key K) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if ele, exists := l.data[key]; exists {
		l.removeElement(ele)
		return nil
	}
	return ErrCacheNotFound
}

// LoadAndDelete 从 LRUCache 中删除与指定键关联的值，并返回该值。
// 如果键不存在，则返回ErrCacheNotFound错误。
func (l *LRUCache[K, V]) LoadAndDelete(ctx context.Context, key K) (Item[K, V], error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var zeroVal Item[K, V]
	if ele, exists := l.data[key]; exists {
		l.removeElement(ele)
		return ele.Value.(Item[K, V]), nil
	}
	return zeroVal, ErrCacheNotFound
}

func (l *LRUCache[K, V]) startGC() {
	if l.gcInterval.Milliseconds() == 0 {
		return
	}
	ticker := time.NewTicker(l.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.mu.Lock()
			l.gcRandom()
			l.mu.Unlock()
		}
	}
}

func (l *LRUCache[K, V]) gcRandom() {
	length := l.Len()

	if length == 0 {
		return
	}

	randomIndex := rand.Intn(length)

	for i := 0; i < l.gcRandomDeletionStep && randomIndex < length; i++ {
		ele := l.list.Front()
		for j := 0; j < randomIndex; j++ {
			ele = ele.Next()
		}

		if ele != nil {
			expiration := ele.Value.(Item[K, V]).Expiration
			if !expiration.IsZero() && time.Now().After(expiration) {
				l.removeElement(ele)
			}
		}

		randomIndex++
	}
}

// Len 返回 LRUCache 中的元素个数。
func (l *LRUCache[K, V]) Len() int {
	return len(l.data)
}

// removeOldest 移除最近最少使用的元素。
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
	l.mu.RLock()
	defer l.mu.RUnlock()

	ele := l.list.Front()
	for ele != nil {
		fmt.Println(ele.Value.(V))
		ele = ele.Next()
	}
}
