package main

import "fmt"

func main() {
	input := []int{1, 3, 5, 6}
	output := searchInsert(input, 0)
	fmt.Println(output)
}

func searchInsert(nums []int, target int) int {
	n := len(nums)
	left := 0
	right := n - 1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			return mid
		}
	}
	return left
}
