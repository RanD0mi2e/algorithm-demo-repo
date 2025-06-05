package main

import "fmt"

func combine(n int, k int) [][]int {
	var result [][]int
	var current []int

	var backTrace func(start int)
	backTrace = func(start int) {
		if len(current) == k {
			combination := make([]int, k)
			copy(combination, current)
			result = append(result, combination)
			return
		}

		for i := start; i <= n; i++ {
			current = append(current, i)
			backTrace(i + 1)
			current = current[:len(current)-1]
		}
	}

	backTrace(1)
	return result
}

func main() {
	result := combine(4, 2)
	for i, combination := range result {
		fmt.Printf("组合%d: %v\n", i+1, combination)
	}

}
