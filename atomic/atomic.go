package atomic

import "sync/atomic"

// Value 原子值
type Value[T any] struct {
	value atomic.Value
}

// NewAtomic 创建一个原子值
func NewAtomic[T any](value T) *Value[T] {
	v := atomic.Value{}
	v.Store(value)
	return &Value[T]{value: v}
}

// Load 加载值
func (v *Value[T]) Load() T {
	return v.value.Load().(T)
}

// Store 存储值
func (v *Value[T]) Store(value T) {
	v.value.Store(value)
}

// Swap 交换值，并返回旧值
func (v *Value[T]) Swap(value T) T {
	return v.value.Swap(value).(T)
}

// CompareAndSwap 如果当前值等于 old，则将值设置为 new，并返回 true，否则返回 false
func (v *Value[T]) CompareAndSwap(old, new T) bool {
	return v.value.CompareAndSwap(old, new)
}
