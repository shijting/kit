package syncx

import "sync"

type Map[K comparable, V any] struct {
	m sync.Map
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{m: sync.Map{}}
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return value, false
	}
	return v.(V), true
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var v any
	v, loaded = m.m.LoadOrStore(key, value)
	if v != nil {
		actual = v.(V)
	}
	return actual, loaded
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var v any
	v, loaded = m.m.LoadAndDelete(key)
	if v != nil {
		value = v.(V)
	}
	return value, loaded

}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		var k K
		var v V
		if value != nil {
			v = value.(V)
		}
		if key != nil {
			k = key.(K)
		}
		return f(k, v)
	})
}

func (m *Map[K, V]) Len() int64 {
	var cnt int64
	m.m.Range(func(key, value any) bool {
		cnt++
		return true
	})
	return cnt
}
