package atomic

import "sync/atomic"

type Value[T any] struct {
	value atomic.Value
}

func NewAtomic[T any](value T) *Value[T] {
	v := atomic.Value{}
	v.Store(value)
	return &Value[T]{value: v}
}

func (v *Value[T]) Load() T {
	return v.value.Load().(T)
}

func (v *Value[T]) Store(value T) {
	v.value.Store(value)
}

func (v *Value[T]) Swap(value T) T {
	return v.value.Swap(value).(T)
}

func (v *Value[T]) CompareAndSwap(old, new T) bool {
	return v.value.CompareAndSwap(old, new)
}
