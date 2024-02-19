package cache

import (
	"context"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (Item[K, V], bool)
	Set(ctx context.Context, key K, val V, expiration time.Duration) error
	Delete(ctx context.Context, key K) error

	LoadAndDelete(ctx context.Context, key K) (Item[K, V], error)
}
