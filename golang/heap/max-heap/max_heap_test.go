package maxheap

import (
	"testing"
)

func TestNewMaxheap(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	if maxHeap == nil {
		t.Error("NewMaxheap should return a non-nil pointer")
	}

	if maxHeap.Size() != 0 {
		t.Errorf("New heap should have size 0, got %d", maxHeap.Size())
	}
}

func TestInsertAndSize(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 测试插入单个元素
	maxHeap.Insert(10)
	if maxHeap.Size() != 1 {
		t.Errorf("Expected size 1, got %d", maxHeap.Size())
	}

	// 测试插入多个元素
	maxHeap.Insert(20)
	maxHeap.Insert(5)
	maxHeap.Insert(15)
	if maxHeap.Size() != 4 {
		t.Errorf("Expected size 4, got %d", maxHeap.Size())
	}
}

func TestPeek(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 测试空堆
	val, ok := maxHeap.Peek()
	if ok {
		t.Error("Peek on empty heap should return false")
	}
	if val != -1 {
		t.Errorf("Peek on empty heap should return -1, got %d", val)
	}

	// 测试单个元素
	maxHeap.Insert(10)
	val, ok = maxHeap.Peek()
	if !ok {
		t.Error("Peek should return true for non-empty heap")
	}
	if val != 10 {
		t.Errorf("Expected peek value 10, got %d", val)
	}

	// 测试多个元素，确保最大值在顶部
	maxHeap.Insert(20)
	maxHeap.Insert(5)
	val, ok = maxHeap.Peek()
	if !ok {
		t.Error("Peek should return true for non-empty heap")
	}
	if val != 20 {
		t.Errorf("Expected peek value 20, got %d", val)
	}
}

func TestExtractMax(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 测试空堆
	val, ok := maxHeap.ExtractMax()
	if ok {
		t.Error("ExtractMax on empty heap should return false")
	}
	if val != -1 {
		t.Errorf("ExtractMax on empty heap should return -1, got %d", val)
	}

	// 测试单个元素
	maxHeap.Insert(10)
	val, ok = maxHeap.ExtractMax()
	if !ok {
		t.Error("ExtractMax should return true for non-empty heap")
	}
	if val != 10 {
		t.Errorf("Expected extracted value 10, got %d", val)
	}
	if maxHeap.Size() != 0 {
		t.Errorf("Heap should be empty after extracting last element, size: %d", maxHeap.Size())
	}

	// 测试多个元素
	maxHeap.Insert(10)
	maxHeap.Insert(20)
	maxHeap.Insert(5)
	maxHeap.Insert(15)

	// 应该按降序提取
	expected := []int{20, 15, 10, 5}
	for i, expectedVal := range expected {
		val, ok := maxHeap.ExtractMax()
		if !ok {
			t.Errorf("ExtractMax should return true, iteration %d", i)
		}
		if val != expectedVal {
			t.Errorf("Expected %d, got %d at iteration %d", expectedVal, val, i)
		}
	}

	if maxHeap.Size() != 0 {
		t.Errorf("Heap should be empty after extracting all elements, size: %d", maxHeap.Size())
	}
}

func TestHeapProperty(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 插入一系列数字并验证堆属性
	values := []int{3, 7, 4, 10, 12, 9, 6, 15, 16, 4, 6, 19, 1}
	for _, val := range values {
		maxHeap.Insert(val)
	}

	// 提取所有元素，应该是降序
	var extracted []int
	for maxHeap.Size() > 0 {
		val, ok := maxHeap.ExtractMax()
		if !ok {
			t.Error("ExtractMax should return true for non-empty heap")
		}
		extracted = append(extracted, val)
	}

	// 验证是否为降序
	for i := 1; i < len(extracted); i++ {
		if extracted[i-1] < extracted[i] {
			t.Errorf("Heap property violated: %d should be >= %d", extracted[i-1], extracted[i])
		}
	}
}

func TestDuplicateValues(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 测试重复值
	maxHeap.Insert(10)
	maxHeap.Insert(10)
	maxHeap.Insert(10)

	if maxHeap.Size() != 3 {
		t.Errorf("Expected size 3, got %d", maxHeap.Size())
	}

	for i := 0; i < 3; i++ {
		val, ok := maxHeap.ExtractMax()
		if !ok {
			t.Errorf("ExtractMax should return true, iteration %d", i)
		}
		if val != 10 {
			t.Errorf("Expected 10, got %d at iteration %d", val, i)
		}
	}
}

func TestNegativeValues(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 测试负数
	maxHeap.Insert(-5)
	maxHeap.Insert(-10)
	maxHeap.Insert(-1)

	val, ok := maxHeap.Peek()
	if !ok {
		t.Error("Peek should return true for non-empty heap")
	}
	if val != -1 {
		t.Errorf("Expected peek value -1, got %d", val)
	}

	// 提取应该是 -1, -5, -10
	expected := []int{-1, -5, -10}
	for i, expectedVal := range expected {
		val, ok := maxHeap.ExtractMax()
		if !ok {
			t.Errorf("ExtractMax should return true, iteration %d", i)
		}
		if val != expectedVal {
			t.Errorf("Expected %d, got %d at iteration %d", expectedVal, val, i)
		}
	}
}

func TestMixedValues(t *testing.T) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 测试正数、负数和零的混合
	values := []int{0, -5, 10, -3, 7, 0, -1}
	for _, val := range values {
		maxHeap.Insert(val)
	}

	// 第一个应该是最大值 10
	val, ok := maxHeap.Peek()
	if !ok {
		t.Error("Peek should return true for non-empty heap")
	}
	if val != 10 {
		t.Errorf("Expected peek value 10, got %d", val)
	}

	// 提取所有元素并验证降序
	var extracted []int
	for maxHeap.Size() > 0 {
		val, ok := maxHeap.ExtractMax()
		if !ok {
			t.Error("ExtractMax should return true for non-empty heap")
		}
		extracted = append(extracted, val)
	}

	for i := 1; i < len(extracted); i++ {
		if extracted[i-1] < extracted[i] {
			t.Errorf("Heap property violated: %d should be >= %d", extracted[i-1], extracted[i])
		}
	}
}

// 基准测试
func BenchmarkInsert(b *testing.B) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maxHeap.Insert(i)
	}
}

func BenchmarkExtractMax(b *testing.B) {
	h := &Maxheap{}
	maxHeap := h.NewMaxheap()

	// 预填充堆
	for i := 0; i < b.N; i++ {
		maxHeap.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maxHeap.ExtractMax()
	}
}
