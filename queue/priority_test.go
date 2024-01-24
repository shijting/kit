package queue

import (
	"testing"
)

//func TestPriorityQueue(t *testing.T) {
//	//fixedSizeHeap := NewFixedSizeHeap[string](5)
//	//
//	//arr := []int{90, 2, 1, 10, 9, 5, 4, 3, 8, 7, 6, 20, 30, 40, 100, 0}
//	//for _, i := range arr {
//	//	item := &Item[string]{Value: fmt.Sprintf("Task %d", i), Priority: i}
//	//	fixedSizeHeap.Push(item)
//	//}
//	//
//	//// 弹出元素
//	//for !fixedSizeHeap.IsEmpty() {
//	//	item := fixedSizeHeap.Pop()
//	//	fmt.Println(item.Value)
//	//}
//
//	type TestUser struct {
//		Name string
//		Age  int
//	}
//
//	fixedSizeHeap := NewFixedSizeHeap[TestUser](5)
//	arr := []int{90, 2, 1, 10, 9, 5, 4, 3, 8, 7, 6, 20, 30, 40, 100, 0}
//	for _, i := range arr {
//		item := &Item[TestUser]{Value: TestUser{
//			Name: fmt.Sprintf("Task %d", i),
//			Age:  i,
//		}, Priority: i}
//		fixedSizeHeap.Push(item)
//	}
//
//	println("fixedSizeHeap.Len() = ", fixedSizeHeap.Len())
//
//	for !fixedSizeHeap.IsEmpty() {
//		item := fixedSizeHeap.Pop()
//		fmt.Println(item.Value)
//	}
//}

func TestFixedSizeHeap(t *testing.T) {
	fixedSizeHeap := NewFixedSizeHeap[int](3)
	fixedSizeHeap.Push(nil)

	if fixedSizeHeap.Len() != 0 {
		t.Errorf("Expected length 0, got %d", fixedSizeHeap.Len())
	}

	poppedItem := fixedSizeHeap.Pop()

	if poppedItem != nil {
		t.Errorf("Expected poppedItem to be nil, got %v", poppedItem)
	}

	fixedSizeHeap.Push(&Item[int]{Value: 1, Priority: 1})
	fixedSizeHeap.Push(&Item[int]{Value: 2, Priority: 2})
	fixedSizeHeap.Push(&Item[int]{Value: 3, Priority: 3})

	if fixedSizeHeap.Len() != 3 {
		t.Errorf("Expected length 3, got %d", fixedSizeHeap.Len())
	}

	if fixedSizeHeap.heap[0].Priority != 1 {
		t.Errorf("Expected priority 1, got %d", fixedSizeHeap.heap[0].Priority)
	}

	fixedSizeHeap.Push(&Item[int]{Value: 4, Priority: 4})

	poppedItem = fixedSizeHeap.Pop()
	if poppedItem.Priority != 2 {
		t.Errorf("Expected priority 2, got %d", poppedItem.Priority)
	}

	if fixedSizeHeap.IsEmpty() {
		t.Errorf("Expected IsEmpty to be false, got true")
	}

	poppedItem = fixedSizeHeap.Pop()
	if poppedItem.Priority != 3 {
		t.Errorf("Expected priority 3, got %d", poppedItem.Priority)
	}

	// 再次 Pop
	poppedItem = fixedSizeHeap.Pop()
	if poppedItem.Priority != 4 {
		t.Errorf("Expected priority 4, got %d", poppedItem.Priority)
	}

	// 测试 IsEmpty 方法
	if !fixedSizeHeap.IsEmpty() {
		t.Errorf("Expected IsEmpty to be true, got false")
	}
}
