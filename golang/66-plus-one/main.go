package main

import "fmt"

func main() {
	input := []int{9, 9, 9, 9}
	output := plusOne(input)
	fmt.Printf("result: %v", output)
}

func plusOne(digits []int) []int {
	n := len(digits)
	// case1：数组存在不为9的数字
	for i := n - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i] += 1
			for j := i + 1; j < n; j++ {
				digits[j] = 0
			}
			return digits
		}
	}

	// case2:数组均为9
	digits = make([]int, n+1)
	digits[0] = 1
	return digits
}
