package cache

import (
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
					lru.Set(1, "one", 0)
				},
				func(lru *LRUCache[int, string]) {
					val, ok := lru.Get(1)
					assert.Equal(t, ok, true)
					assert.Equal(t, val.Value, "one")
				},
				func(lru *LRUCache[int, string]) {
					lru.Set(2, "two", 0)
				},
				func(lru *LRUCache[int, string]) {
					val, ok := lru.Get(2)
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
					lru.Set(1, "one", 0)
				},
				func(lru *LRUCache[int, string]) {
					lru.Set(2, "two", 0)
				},
				func(lru *LRUCache[int, string]) {
					lru.Set(3, "three", 0)
				},
				func(lru *LRUCache[int, string]) {
					_, ok := lru.Get(1)

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
					lru.Set(1, "one", time.Second)
					time.Sleep(time.Second)
					_, ok := lru.Get(1)
					assert.Equal(t, ok, false)
				},
				func(lru *LRUCache[int, string]) {
					lru.Set(2, "two", time.Second*2)
					time.Sleep(time.Second * 2)
					_, ok := lru.Get(2)
					assert.Equal(t, ok, false)
				},
				func(lru *LRUCache[int, string]) {
					lru.Set(3, "three", 3*time.Second)
					time.Sleep(time.Second)
					data, ok := lru.Get(3)
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
