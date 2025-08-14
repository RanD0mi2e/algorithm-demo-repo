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

	mf := Mf.Constructor()
	mf.AddNum(3)
	mf.AddNum(2)
	result1 := mf.FindMedian()
	fmt.Printf("result1: %v\n", result1) // expect 1.5
	mf.AddNum(1)
	result2 := mf.FindMedian()
	fmt.Printf("result2: %v\n", result2) // expect 2
}
