package medianfinder

import "container/heap"

type MinHeap []int

func (mh MinHeap) Less(i, j int) bool {
	return mh[i] < mh[j]
}

func (mh MinHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
}

func (mh MinHeap) Len() int {
	return len(mh)
}

func (mh *MinHeap) Pop() any {
	old := *mh
	n := len(old)
	item := old[n-1]
	*mh = old[0 : n-1]
	return item
}

func (mh *MinHeap) Push(x any) {
	*mh = append(*mh, x.(int))
}

func (mh MinHeap) Peek() int {
	return mh[0]
}

type MaxHeap []int

func (mh MaxHeap) Less(i, j int) bool {
	return mh[i] > mh[j]
}

func (mh MaxHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
}

func (mh MaxHeap) Len() int {
	return len(mh)
}

func (mh *MaxHeap) Pop() any {
	old := *mh
	n := len(old)
	item := old[n-1]
	*mh = old[0 : n-1]
	return item
}

func (mh *MaxHeap) Push(x any) {
	*mh = append(*mh, x.(int))
}

func (mh MaxHeap) Peek() int {
	return mh[0]
}

type MedianFinder struct {
	MinHeap *MinHeap
	MaxHeap *MaxHeap
}

func Constructor() MedianFinder {
	minHp := &MinHeap{}
	maxHp := &MaxHeap{}
	heap.Init(minHp)
	heap.Init(maxHp)

	return MedianFinder{
		MinHeap: minHp,
		MaxHeap: maxHp,
	}
}

func (this *MedianFinder) AddNum(num int) {
	if len(*this.MaxHeap) == 0 || num < this.MaxHeap.Peek() {
		heap.Push(this.MaxHeap, num)
	} else {
		heap.Push(this.MinHeap, num)
	}
	this.Balance()
}

func (this *MedianFinder) Balance() {
	minHeapLen, maxHeapLen := len(*this.MinHeap), len(*this.MaxHeap)
	if maxHeapLen > minHeapLen+1 {
		val := heap.Pop(this.MaxHeap).(int)
		heap.Push(this.MinHeap, val)
	} else if minHeapLen > maxHeapLen+1 {
		val := heap.Pop(this.MinHeap).(int)
		heap.Push(this.MaxHeap, val)
	}
}

func (this *MedianFinder) FindMedian() float64 {
	minL, maxL := len(*this.MinHeap), len(*this.MaxHeap)
	if minL == maxL {
		return (float64(this.MinHeap.Peek()) + float64(this.MaxHeap.Peek())) / 2
	} else if minL > maxL {
		return float64(this.MinHeap.Peek())
	} else {
		return float64(this.MaxHeap.Peek())
	}
}
