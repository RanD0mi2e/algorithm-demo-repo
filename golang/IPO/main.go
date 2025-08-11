package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type Heap struct {
	sort.IntSlice
}

func (h *Heap) Less(i, j int) bool {
	return h.IntSlice[i] > h.IntSlice[j]
}

func (h *Heap) Push(val any) {
	h.IntSlice = append(h.IntSlice, val.(int))
}

func (h *Heap) Pop() any {
	arr := h.IntSlice
	val := arr[len(arr)-1]
	h.IntSlice = arr[:len(arr)-1]
	return val
}

func main() {
	result := findMaximizedCapital(4, 0, []int{1, 1, 3}, []int{0, 0, 3})
	fmt.Printf("result: %v", result)
}

func findMaximizedCapital(k int, w int, profits []int, capital []int) int {
	n := len(profits)
	type pair struct {
		capital int
		profit  int
	}

	arr := make([]pair, n)
	// o(n)
	for i, p := range profits {
		arr[i] = pair{
			capital: capital[i],
			profit:  p,
		}
	}
	// o(nlogn)
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].capital < arr[j].capital
	})

	hp := &Heap{}
	// o(nlogn + klogn)
	for curr := 0; k > 0; k-- {
		for curr < n && arr[curr].capital <= w {
			heap.Push(hp, arr[curr].profit)
			curr++
		}

		if hp.Len() == 0 {
			break
		}

		w += heap.Pop(hp).(int)
	}

	return w
}
