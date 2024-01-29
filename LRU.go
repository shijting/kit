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

func (l *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	if v, ok := l.data.Load(key); ok {
		// 移动到最前面
		l.list.MoveToFront(v)
		return v.Value.(V), true
	}
	var v V
	return v, false
}

func (l *LRUCache[K, V]) Set(key K, value V) {
	if v, ok := l.data.Load(key); ok {
		v.Value = value
		// 移动到最前面
		l.list.MoveToFront(v)
	} else {
		// 新增 并且移动到最前面
		l.data.Store(key, l.list.PushFront(value))
	}
}

func (l *LRUCache[K, V]) Print() {
	ele := l.list.Front()
	for ele != nil {
		fmt.Println(ele.Value.(V))
		ele = ele.Next()
	}
}
