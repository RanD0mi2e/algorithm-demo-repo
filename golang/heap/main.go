package main

import (
	topk "algorithm/heap/top-k"
	"fmt"
)

func main() {
	// input := []int{3, 1, 8, 5, 2}
	// heapsort.Heapsort(input)
	// fmt.Printf("排序后的数组: %v\n", input)
	input := []int{1, 2, 3, 4, 5}
	output := topk.TopK(input, 3)
	fmt.Printf("result: %v", output)
}
