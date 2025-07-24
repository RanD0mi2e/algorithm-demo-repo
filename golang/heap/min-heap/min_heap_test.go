package Minheap

import (
	"testing"
)

func TestMinheap(t *testing.T) {
	// 测试创建空堆
	h := NewMinheap()
	if h.Size() != 0 {
		t.Errorf("Expected size 0, got %d", h.Size())
	}

	// 测试空堆的Peek和ExtractMin
	if _, ok := h.Peek(); ok {
		t.Error("Expected Peek to return false for empty heap")
	}
	if _, ok := h.ExtractMin(); ok {
		t.Error("Expected ExtractMin to return false for empty heap")
	}

	// 测试插入单个元素
	h.Insert(5)
	if h.Size() != 1 {
		t.Errorf("Expected size 1, got %d", h.Size())
	}
	if val, ok := h.Peek(); !ok || val != 5 {
		t.Errorf("Expected Peek to return 5, got %d", val)
	}

	// 测试插入多个元素并验证最小堆性质
	elements := []int{3, 8, 1, 9, 2, 7}
	for _, elem := range elements {
		h.Insert(elem)
	}

	// 验证堆顶始终是最小值
	if val, ok := h.Peek(); !ok || val != 1 {
		t.Errorf("Expected min value 1, got %d", val)
	}

	// 测试ExtractMin - 应该按升序返回所有元素
	expected := []int{1, 2, 3, 5, 7, 8, 9}
	for i, expectedVal := range expected {
		if val, ok := h.ExtractMin(); !ok || val != expectedVal {
			t.Errorf("ExtractMin %d: expected %d, got %d", i, expectedVal, val)
		}
	}

	// 验证堆为空
	if h.Size() != 0 {
		t.Errorf("Expected empty heap, size is %d", h.Size())
	}
}

func TestMinheapEdgeCases(t *testing.T) {
	h := NewMinheap()

	// 测试重复元素
	h.Insert(5)
	h.Insert(5)
	h.Insert(3)
	h.Insert(3)

	if val, _ := h.ExtractMin(); val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}
	if val, _ := h.ExtractMin(); val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}
	if val, _ := h.ExtractMin(); val != 5 {
		t.Errorf("Expected 5, got %d", val)
	}
	if val, _ := h.ExtractMin(); val != 5 {
		t.Errorf("Expected 5, got %d", val)
	}

	// 测试负数
	h2 := NewMinheap()
	h2.Insert(-1)
	h2.Insert(-5)
	h2.Insert(0)
	h2.Insert(3)

	if val, _ := h2.Peek(); val != -5 {
		t.Errorf("Expected -5, got %d", val)
	}
}

func TestMinheapLargeDataset(t *testing.T) {
	h := NewMinheap()

	// 插入大量数据测试性能和正确性
	data := []int{100, 50, 150, 25, 75, 125, 175, 10, 30, 60, 80}
	for _, val := range data {
		h.Insert(val)
	}

	// 验证提取的顺序是否正确
	prev := -1
	for h.Size() > 0 {
		val, _ := h.ExtractMin()
		if prev != -1 && val < prev {
			t.Errorf("Heap property violated: %d < %d", val, prev)
		}
		prev = val
	}
}
