package kit

import (
	"container/list"
	"fmt"
	"github.com/shijting/kit/syncx"
)

type LRUCache[K comparable, V any] struct {
	list *list.List
	data *syncx.Map[K, *list.Element]
}

func NewLRUCache[K comparable, V any]() *LRUCache[K, V] {
	return &LRUCache[K, V]{
		list: list.New(),
		data: syncx.NewMap[K, *list.Element](),
	}
}

type Data[K comparable, V any] struct {
	Key   K
	Value V
	//	TODO 过期时间
}

func NewData[K comparable, V any](key K, value V) Data[K, V] {
	return Data[K, V]{Key: key, Value: value}
}

func (l *LRUCache[K, V]) Get(key K) (value Data[K, V], ok bool) {
	if v, ok := l.data.Load(key); ok {
		// 移动到最前面
		l.list.MoveToFront(v)
		return v.Value.(Data[K, V]), true
	}
	var v Data[K, V]
	return v, false
}

func (l *LRUCache[K, V]) Set(key K, value Data[K, V]) {
	if v, ok := l.data.Load(key); ok {
		v.Value = value
		// 移动到最前面
		l.list.MoveToFront(v)
	} else {
		// 新增 并且移动到最前面
		l.data.Store(key, l.list.PushFront(value))
	}
}

// RemoveBack 删除最后一个
func (l *LRUCache[K, V]) RemoveBack() {
	ele := l.list.Back()
	if ele != nil {
		l.list.Remove(ele)
		key := ele.Value.(Data[K, V]).Key
		l.data.LoadAndDelete(key)
	}
}

func (l *LRUCache[K, V]) Print() {
	ele := l.list.Front()
	for ele != nil {
		fmt.Println(ele.Value.(V))
		ele = ele.Next()
	}
}
