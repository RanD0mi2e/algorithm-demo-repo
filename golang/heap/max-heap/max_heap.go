package maxheap

type Maxheap struct {
	heap []int
}

func (h *Maxheap) NewMaxheap() *Maxheap {
	return &Maxheap{
		heap: make([]int, 0),
	}
}

func (h *Maxheap) leftChild(i int) int {
	return 2*i + 1
}

func (h *Maxheap) rightChild(i int) int {
	return 2*i + 2
}

func (h *Maxheap) parent(i int) int {
	return (i - 1) / 2
}

func (h *Maxheap) swin(i int) {
	for i > 0 && h.heap[i] > h.heap[h.parent(i)] {
		h.heap[i], h.heap[h.parent(i)] = h.heap[h.parent(i)], h.heap[i]
		i = h.parent(i)
	}
}

func (h *Maxheap) sink(i int) {
	for h.leftChild(i) < len(h.heap) {
		maxChild := h.leftChild(i)
		if h.rightChild(i) < len(h.heap) && h.heap[h.rightChild(i)] > h.heap[h.leftChild(i)] {
			maxChild = h.rightChild(i)
		}
		if h.heap[i] >= h.heap[maxChild] {
			break
		}
		h.heap[i], h.heap[maxChild] = h.heap[maxChild], h.heap[i]
		i = maxChild
	}
}

func (h *Maxheap) Size() int {
	return len(h.heap)
}

func (h *Maxheap) Insert(val int) {
	h.heap = append(h.heap, val)
	h.swin(h.Size() - 1)
}

func (h *Maxheap) ExtractMax() (int, bool) {
	if h.Size() == 0 {
		return -1, false
	}
	max := h.heap[0]
	h.heap[0] = h.heap[h.Size()-1]
	h.heap = h.heap[:h.Size()-1]
	if h.Size() > 0 {
		h.sink(0)
	}
	return max, true
}

func (h *Maxheap) Peek() (int, bool) {
	if h.Size() == 0 {
		return -1, false
	}
	return h.heap[0], true
}

// heap堆外部api
// NewMaxheap创建最大堆 x
// Insert插入元素 x
// ExtractMax获取最大值 x
// Peek查看最大值 x
// Size查看当前堆大小 x

// heap堆内部api
// sink小元素下沉 x
// swim大元素上浮 x
// leftChild根据父元素索引获取左子节点 x
// rightChild根据父元素索引获取右子节点 x
// parent获取当前节点的父元素索引 x
