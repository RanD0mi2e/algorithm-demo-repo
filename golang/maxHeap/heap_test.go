package maxHeap

import (
	"testing"
)

func TestMaxHeap_Insert(t *testing.T) {
	hp := NewMaxHeap(10)
	hp.Insert(5)
	hp.Insert(10)
	hp.Insert(3)

	if hp.Pq[1] != 10 {
		t.Errorf("优先队列第一个元素预期为 10, 但实际为 %d", hp.Pq[1])
	}
	if hp.Pq[2] != 5 {
		t.Errorf("优先队列第二个元素预期为 5, 但实际为 %d", hp.Pq[2])
	}
	if hp.Pq[3] != 3 {
		t.Errorf("优先队列第三个元素预期为 3, 但实际为 %d", hp.Pq[3])
	}
}

func TestMaxHeap_DeleteMax(t *testing.T) {
	heap := NewMaxHeap(10)
	heap.Insert(5)
	heap.Insert(10)
	heap.Insert(3)

	max := heap.DeleteMax()
	if max != 10 {
		t.Errorf("Expected max element to be 10, got %d", max)
	}
	if heap.Pq[1] != 5 {
		t.Errorf("Expected new max element to be 5, got %d", heap.Pq[1])
	}
	if heap.N != 2 {
		t.Errorf("Expected heap size to be 2, got %d", heap.N)
	}
}

func TestMaxHeap_IsEmpty(t *testing.T) {
	heap := NewMaxHeap(10)
	if !heap.IsEmpty() {
		t.Errorf("Expected heap to be empty")
	}
	heap.Insert(5)
	if heap.IsEmpty() {
		t.Errorf("Expected heap to not be empty")
	}
}

func TestMaxHeap_Size(t *testing.T) {
	heap := NewMaxHeap(10)
	if heap.Size() != 0 {
		t.Errorf("Expected heap size to be 0, got %d", heap.Size())
	}
	heap.Insert(5)
	if heap.Size() != 1 {
		t.Errorf("Expected heap size to be 1, got %d", heap.Size())
	}
}
