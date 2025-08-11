package medianfinder

import "container/heap"

type MinHeap []int

func (h *MinHeap) Len() int { return len(*h) }

func (h *MinHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *MinHeap) Pop() interface{} {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}

func (h *MinHeap) Push(val interface{}) {
	*h = append(*h, val.(int))
}

func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeap) Peek() int {
	return (*h)[0]
}

type MaxHeap []int

func (h *MaxHeap) Len() int { return len(*h) }

func (h *MaxHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *MaxHeap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}

func (h *MaxHeap) Push(val any) {
	*h = append(*h, val.(int))
}

func (h *MaxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MaxHeap) Peek() int {
	return (*h)[0]
}

type MedianFinder struct {
	minheap *MinHeap
	maxheap *MaxHeap
}

func NewMedianFinder() *MedianFinder {
	minheap := &MinHeap{}
	maxheap := &MaxHeap{}
	heap.Init(minheap)
	heap.Init(maxheap)
	return &MedianFinder{
		minheap: minheap,
		maxheap: maxheap,
	}
}

func (mf *MedianFinder) Balance() {
	minLen, maxLen := len(*mf.minheap), len(*mf.maxheap)
	if maxLen > minLen+1 {
		val := heap.Pop(mf.maxheap).(int)
		heap.Push(mf.minheap, val)
	} else if minLen > maxLen+1 {
		val := heap.Pop(mf.minheap).(int)
		heap.Push(mf.maxheap, val)
	}
}

func (mf *MedianFinder) AddNumber(val int) {
	maxLen := len(*mf.maxheap)
	if maxLen == 0 || val <= (*mf.maxheap)[0] {
		heap.Push(mf.maxheap, val)
	} else {
		heap.Push(mf.minheap, val)
	}
	mf.Balance()
}

func (mf *MedianFinder) GetMedian() float64 {
	var result float64
	minLen, maxLen := (*mf.minheap).Len(), (*mf.maxheap).Len()
	if minLen == maxLen {
		result = float64((*mf.minheap).Peek()+(*mf.maxheap).Peek()) / 2.0
	} else if minLen > maxLen {
		result = float64((*mf.minheap).Peek())
	} else {
		result = float64((*mf.maxheap).Peek())
	}
	return result
}
