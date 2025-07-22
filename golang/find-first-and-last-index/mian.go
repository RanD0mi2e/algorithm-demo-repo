package main

import "fmt"

func main() {
	input := []int{1}
	output := searchRange(input, 1)
	fmt.Printf("result: %+v", output)
}

func searchRange(nums []int, target int) []int {
	length := len(nums)
	if length == 0 {
		return []int{-1, -1}
	}

	findBound := func(boundType string) int {
		left, right := 0, length-1
		result := -1
		for left <= right {
			mid := left + (right-left)>>1
			if target < nums[mid] {
				right = mid - 1
			} else if target > nums[mid] {
				left = mid + 1
			} else {
				result = mid
				if boundType == "left" {
					right = mid - 1
				}
				if boundType == "right" {
					left = mid + 1
				}
			}
		}
		return result
	}

	firstPos := findBound("left")
	if firstPos == -1 {
		return []int{-1, -1}
	}
	lastPos := findBound("right")

	return []int{firstPos, lastPos}
}
