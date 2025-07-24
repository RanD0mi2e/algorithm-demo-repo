package topk

func TopK(arr []int, k int) []int {
	hp := newMaxheap()
	result := []int{}
	for _, num := range arr {
		if hp.size() < k {
			hp.insert(num)
		} else {
			if top, ok := hp.peek(); ok && num < top {
				hp.extract()
				hp.insert(num)
			}
		}
	}
	result = append(result, hp.heap...)
	return result
}

type maxheap struct {
	heap []int
}

func newMaxheap() *maxheap {
	return &maxheap{
		heap: make([]int, 0),
	}
}

func (h *maxheap) parent(i int) int {
	return (i - 1) / 2
}

func (h *maxheap) left(i int) int {
	return 2*i + 1
}

func (h *maxheap) right(i int) int {
	return 2*i + 2
}

func (h *maxheap) swim(i int) {
	for i > 0 && h.heap[h.parent(i)] < h.heap[i] {
		h.heap[i], h.heap[h.parent(i)] = h.heap[h.parent(i)], h.heap[i]
		i = h.parent(i)
	}
}

func (h *maxheap) sink(i int) {
	for h.heap[h.left(i)] < len(h.heap) {
		childIdx := h.left(i)
		if h.heap[h.right(i)] < len(h.heap) && h.heap[h.right(i)] > h.heap[childIdx] {
			childIdx = h.right(i)
		}
		if h.heap[i] >= h.heap[childIdx] {
			break
		}
		h.heap[i], h.heap[childIdx] = h.heap[childIdx], h.heap[i]
		i = childIdx
	}
}

func (h *maxheap) insert(val int) {
	h.heap = append(h.heap, val)
	h.swim(len(h.heap) - 1)
}

func (h *maxheap) extract() (int, bool) {
	if len(h.heap) == 0 {
		return -1, false
	}
	max := h.heap[0]
	h.heap[0] = h.heap[len(h.heap)-1]
	h.heap = h.heap[:len(h.heap)-1]
	if len(h.heap) > 0 {
		h.sink(0)
	}
	return max, true
}

func (h *maxheap) peek() (int, bool) {
	if len(h.heap) == 0 {
		return -1, false
	}
	return h.heap[0], true
}

func (h *maxheap) size() int {
	return len(h.heap)
}
