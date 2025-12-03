package main

import (
	"fmt"
	"math"
)

func main() {
	input := []int{2, 2, 1, 1, 1, 2, 2}
	output := majorityElement(input)
	fmt.Println(output)
}

// 输入：nums = [2,2,1,1,1,2,2]
// 输出：2
func majorityElement(nums []int) int {
	count := 0
	currNum := math.MinInt
	for _, num := range nums {
		if count == 0 {
			currNum = num
		}
		if num == currNum {
			count++
		} else {
			count--
		}
	}
	return currNum
}
