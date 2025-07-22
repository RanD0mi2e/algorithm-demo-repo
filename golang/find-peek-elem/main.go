package main

import "fmt"

func main() {
	input := []int{5, 4, 2, 3}
	output := findPeakElement(input)
	fmt.Printf("peek elem: %d", output)
}

func findPeakElement(nums []int) int {
	n := len(nums)
	left, right := 0, n-1
	for left < right {
		mid := left + (right-left)>>1
		if nums[mid] > nums[mid+1] {
			right = mid
		} else if nums[mid] < nums[mid+1] {
			left = mid + 1
		}
	}
	return left
}
