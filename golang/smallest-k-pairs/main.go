package main

import "container/heap"

type pair struct {
	sum  int
	i, j int
}

type pairHeap []pair

func (p *pairHeap) Push(val interface{}) {
	*p = append(*p, val.(pair))
}

func (p *pairHeap) Pop() interface{} {
	n := len(*p)
	result := (*p)[n-1]
	*p = (*p)[:n-1]
	return result
}

func (p pairHeap) Less(i, j int) bool {
	return p[i].sum < p[j].sum
}

func (p pairHeap) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p pairHeap) Len() int {
	return len(p)
}

func main() {}

func kSmallestPairs(nums1 []int, nums2 []int, k int) [][]int {
	if len(nums1) == 0 || len(nums2) == 0 || k == 0 {
		return [][]int{}
	}

	minHeap := &pairHeap{}
	heap.Init(minHeap)

	for i := 0; i < len(nums1) && i < k; i++ {
		heap.Push(minHeap, pair{
			sum: nums1[i] + nums2[0],
			i:   i,
			j:   0,
		})
	}

	result := make([][]int, 0, k)
	for len(result) < k && len(*minHeap) > 0 {
		curPair := heap.Pop(minHeap).(pair)
		result = append(result, []int{nums1[curPair.i], nums2[curPair.j]})
		if curPair.j+1 < len(nums2) {
			heap.Push(minHeap, pair{
				sum: nums1[curPair.i] + nums2[curPair.j+1],
				i:   curPair.i,
				j:   curPair.j + 1,
			})
		}
	}

	return result
}
