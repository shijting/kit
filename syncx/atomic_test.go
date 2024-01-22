package syncx

import (
	"testing"
)

func TestAtomic(t *testing.T) {
	originalValue := 42
	atomicValue := NewAtomic(originalValue)

	loadedValue := atomicValue.Load()
	if loadedValue != originalValue {
		t.Errorf("Expected %v, got %v", originalValue, loadedValue)
	}

	newValue := 100
	atomicValue.Store(newValue)
	loadedValue = atomicValue.Load()
	if loadedValue != newValue {
		t.Errorf("Expected %v, got %v", newValue, loadedValue)
	}

	oldValue := atomicValue.Swap(200)
	if oldValue != newValue {
		t.Errorf("Expected %v, got %v", newValue, oldValue)
	}
	loadedValue = atomicValue.Load()
	if loadedValue != 200 {
		t.Errorf("Expected %v, got %v", 200, loadedValue)
	}

	success := atomicValue.CompareAndSwap(200, 300)
	if !success {
		t.Errorf("Expected CompareAndSwap to return true")
	}
	loadedValue = atomicValue.Load()
	if loadedValue != 300 {
		t.Errorf("Expected %v, got %v", 300, loadedValue)
	}

	success = atomicValue.CompareAndSwap(400, 500)
	if success {
		t.Errorf("Expected CompareAndSwap to return false")
	}
	loadedValue = atomicValue.Load()
	if loadedValue != 300 {
		t.Errorf("Expected %v, got %v", 300, loadedValue)
	}
}
