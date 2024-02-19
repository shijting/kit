package cache

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	tests := []struct {
		name         string
		maxCapacity  int
		operations   []func(*LRUCache[int, string]) // Functions to perform on LRUCache
		expectedData map[int]Item[int, string]      // Expected data in LRUCache after operations
	}{
		{
			name:        "Basic Get and Set",
			maxCapacity: 2,
			operations: []func(*LRUCache[int, string]){
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 1, "one", 0)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					val, ok := lru.Get(ctx, 1)
					assert.Equal(t, ok, true)
					assert.Equal(t, val.Value, "one")
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 2, "two", 0)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					val, ok := lru.Get(ctx, 2)
					assert.Equal(t, ok, true)
					assert.Equal(t, val.Value, "two")

				},
			},
			expectedData: map[int]Item[int, string]{
				1: {Key: 1, Value: "one"},
				2: {Key: 2, Value: "two"},
			},
		},
		{
			name:        "LRU Eviction",
			maxCapacity: 2,
			operations: []func(*LRUCache[int, string]){
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 1, "one", 0)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 2, "two", 0)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 3, "three", 0)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					_, ok := lru.Get(ctx, 1)

					assert.Equal(t, ok, false)
				},
			},
			expectedData: map[int]Item[int, string]{
				2: {Key: 2, Value: "two"},
				3: {Key: 3, Value: "three"},
			},
		},
		{
			name:        "Expiration",
			maxCapacity: 2,
			operations: []func(*LRUCache[int, string]){
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 1, "one", time.Second)
					time.Sleep(time.Second)
					_, ok := lru.Get(ctx, 1)
					assert.Equal(t, ok, false)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 2, "two", time.Second*2)
					time.Sleep(time.Second * 2)
					_, ok := lru.Get(ctx, 2)
					assert.Equal(t, ok, false)
				},
				func(lru *LRUCache[int, string]) {
					ctx := context.Background()
					lru.Set(ctx, 3, "three", 3*time.Second)
					time.Sleep(time.Second)
					data, ok := lru.Get(ctx, 3)
					assert.Equal(t, ok, true)
					require.WithinDuration(t, data.Expiration, time.Now().Add(3*time.Second), time.Second*2)
				},
			},
			expectedData: map[int]Item[int, string]{
				3: {Key: 3, Value: "three", Expiration: time.Now().Add(time.Second * 3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lru := NewLRUCache[int, string](tt.maxCapacity)

			for _, op := range tt.operations {
				op(lru)
			}

			// Verify the LRUCache contents
			for key, expectedItem := range tt.expectedData {
				ele, exists := lru.data[key]
				assert.Equal(t, exists, true)

				data := ele.Value.(Item[int, string])
				assert.Equal(t, expectedItem.Key, data.Key)
				assert.Equal(t, expectedItem.Value, data.Value)
				require.WithinDuration(t, expectedItem.Expiration, data.Expiration, time.Second*4)
			}
		})
	}
}

func TestSetAndGet(t *testing.T) {
	cache := NewLRUCache[int, string](3)

	cache.Set(context.Background(), 1, "one", 0)
	cache.Set(context.Background(), 2, "two", 0)
	cache.Set(context.Background(), 3, "three", 0)

	// Test normal retrieval
	val, ok := cache.Get(context.Background(), 1)
	require.Equal(t, ok, true)
	require.Equal(t, val.Value, "one")

	// Test retrieving non-existent key
	val, ok = cache.Get(context.Background(), 4)
	require.Equal(t, ok, false)
	require.Equal(t, val.Value, "")

	// Test updating value
	cache.Set(context.Background(), 1, "updated one", 0)
	val, _ = cache.Get(context.Background(), 1)

	assert.Equal(t, val.Value, "updated one")

	// Test eviction
	cache.Set(context.Background(), 4, "four", 0)
	val, ok = cache.Get(context.Background(), 2)
	require.Equal(t, ok, false)
	require.Equal(t, val.Value, "")
}

func TestDelete(t *testing.T) {
	cache := NewLRUCache[int, string](3)

	cache.Set(context.Background(), 1, "one", 0)
	cache.Set(context.Background(), 2, "two", 0)

	// Test deleting existing key
	err := cache.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, ok := cache.Get(context.Background(), 1)
	if ok {
		t.Errorf("Expected key 1 to be deleted")
	}

	// Test deleting non-existent key
	err = cache.Delete(context.Background(), 3)
	if !errors.Is(err, ErrCacheNotFound) {
		t.Errorf("Expected ErrCacheNotFound, got: %v", err)
	}
}

func TestLoadAndDelete(t *testing.T) {
	cache := NewLRUCache[int, string](3)

	cache.Set(context.Background(), 1, "one", 0)

	// Test loading and deleting existing key
	val, err := cache.LoadAndDelete(context.Background(), 1)
	require.NoError(t, err)
	if val.Value != "one" {
		t.Errorf("Expected value 'one', got %v", val)
	}

	_, ok := cache.Get(context.Background(), 1)
	require.Equal(t, ok, false)

	// Test loading and deleting non-existent key
	_, err = cache.LoadAndDelete(context.Background(), 2)
	assert.Equal(t, err, ErrCacheNotFound)
}

func TestLRUGarbageCollection(t *testing.T) {
	cache := NewLRUCache[int, string](3, WithGCInterval[int, string](100*time.Millisecond))

	cache.Set(context.Background(), 1, "one", 100*time.Millisecond)
	cache.Set(context.Background(), 2, "two", 0)

	// Wait for garbage collection
	time.Sleep(200 * time.Millisecond)

	// Test if expired items are deleted
	_, ok := cache.Get(context.Background(), 1)
	require.Equal(t, ok, false)

	_, ok = cache.Get(context.Background(), 2)
	require.Equal(t, ok, true)
}
