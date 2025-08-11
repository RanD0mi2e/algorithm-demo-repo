package main

import (
	Mf "algorithm/heap/median-finder"
	"fmt"
)

func main() {
	// input := []int{3, 1, 8, 5, 2}
	// heapsort.Heapsort(input)
	// fmt.Printf("排序后的数组: %v\n", input)
	// input := []int{1, 2, 3, 4, 5}
	// output := topk.TopK(input, 3)
	// fmt.Printf("result: %v", output)

	mf := Mf.NewMedianFinder()
	mf.AddNumber(5)
	mf.AddNumber(1)
	mf.AddNumber(3)
	mf.AddNumber(8)
	mf.AddNumber(7)
	mf.AddNumber(0)
	output := mf.GetMedian()
	fmt.Printf("result: %f", output)
}
