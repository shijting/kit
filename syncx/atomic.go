package syncx

import "sync/atomic"

// Atomic 原子值
type Atomic[T any] struct {
	value atomic.Value
}

// NewAtomic 创建一个原子值
func NewAtomic[T any](value T) *Atomic[T] {
	v := atomic.Value{}
	v.Store(value)
	return &Atomic[T]{value: v}
}

// Load 加载值
func (v *Atomic[T]) Load() T {
	return v.value.Load().(T)
}

// Store 存储值
func (v *Atomic[T]) Store(value T) {
	v.value.Store(value)
}

// Swap 交换值，并返回旧值
func (v *Atomic[T]) Swap(value T) T {
	return v.value.Swap(value).(T)
}

// CompareAndSwap 如果当前值等于 old，则将值设置为 new，并返回 true，否则返回 false
func (v *Atomic[T]) CompareAndSwap(old, new T) bool {
	return v.value.CompareAndSwap(old, new)
}
