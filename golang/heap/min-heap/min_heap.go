package Minheap

// heap堆外部api
// NewMinheap创建最小堆 x
// Insert插入元素 x
// ExtractMin获取最小值 x
// Peek查看最小值 x
// Size查看当前堆大小 x

// heap堆内部api
// sink小元素下沉 x
// swim大元素上浮 x
// leftChild根据父元素索引获取左子节点 x
// rightChild根据父元素索引获取右子节点 x
// parent获取当前节点的父元素索引 x

type Minheap struct {
	heap []int
}

func NewMinheap() *Minheap {
	return &Minheap{
		heap: make([]int, 0),
	}
}

func (h *Minheap) leftChild(i int) int {
	return 2*i + 1
}

func (h *Minheap) rightChild(i int) int {
	return 2*i + 2
}

func (h *Minheap) parent(i int) int {
	return (i - 1) / 2
}

func (h *Minheap) swim(i int) {
	for i > 0 && h.heap[i] < h.heap[h.parent(i)] {
		h.heap[i], h.heap[h.parent(i)] = h.heap[h.parent(i)], h.heap[i]
		i = h.parent(i)
	}
}

func (h *Minheap) sink(i int) {
	for h.leftChild(i) < len(h.heap) {
		minChild := h.leftChild(i)
		if h.rightChild(i) < len(h.heap) && h.heap[h.rightChild(i)] < h.heap[h.leftChild(i)] {
			minChild = h.rightChild(i)
		}
		if h.heap[i] <= h.heap[minChild] {
			break
		}
		h.heap[i], h.heap[minChild] = h.heap[minChild], h.heap[i]
		i = minChild
	}
}

func (h *Minheap) Size() int {
	return len(h.heap)
}

func (h *Minheap) Peek() (int, bool) {
	if len(h.heap) == 0 {
		return -1, false
	}
	return h.heap[0], true
}

func (h *Minheap) Insert(val int) {
	h.heap = append(h.heap, val)
	h.swim(len(h.heap) - 1)
}

func (h *Minheap) ExtractMin() (int, bool) {
	if len(h.heap) == 0 {
		return -1, false
	}
	min := h.heap[0]
	h.heap[0] = h.heap[len(h.heap)-1]
	h.heap = h.heap[:len(h.heap)-1]
	if len(h.heap) > 0 {
		h.sink(0)
	}
	return min, true
}
