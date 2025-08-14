package main

import (
	"container/heap"
	"fmt"
	"sort"
)

func main() {
	mf := Constructor()
	mf.AddNum(1)
	mf.AddNum(2)
	r1 := mf.FindMedian()
	fmt.Printf("r1: %v\n", r1)
	mf.AddNum(3)
	r2 := mf.FindMedian()
	fmt.Printf("r2: %v", r2)
}

type hp struct {
	sort.IntSlice
}

func (h *hp) Pop() any {
	old := h.IntSlice
	val := h.IntSlice[len(old)-1]
	h.IntSlice = old[:len(old)-1]
	return val
}

func (h *hp) Push(num any) {
	h.IntSlice = append(h.IntSlice, num.(int))
}

type MedianFinder struct {
	maxQue, minQue hp
}

func Constructor() MedianFinder {
	return MedianFinder{}
}

func (mf *MedianFinder) AddNum(num int) {
	minQ, maxQ := &mf.minQue, &mf.maxQue
	if minQ.Len() == 0 || num < -minQ.IntSlice[0] {
		heap.Push(minQ, -num)
		if minQ.Len() > maxQ.Len()+1 {
			heap.Push(maxQ, -heap.Pop(minQ).(int))
		}
	} else {
		heap.Push(maxQ, num)
		if maxQ.Len() > minQ.Len() {
			heap.Push(minQ, -heap.Pop(maxQ).(int))
		}
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	minQ, maxQ := &mf.minQue, &mf.maxQue
	if minQ.Len() > maxQ.Len() {
		return float64(-minQ.IntSlice[0])
	}
	return float64(maxQ.IntSlice[0]-minQ.IntSlice[0]) / 2
}
