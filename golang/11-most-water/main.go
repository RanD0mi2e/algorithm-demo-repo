package main

import "fmt"

func main() {
	input := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	ouput := maxAreaOptimized(input)
	fmt.Printf("result: %v\n", ouput)
}

// area = height * width
func maxArea(height []int) int {
	left, right, area := 0, len(height)-1, 0
	for left < right {
		width := right - left
		if height[left] < height[right] {
			area = max(area, width*height[left])
			left++
		} else {
			area = max(area, width*height[right])
			right--
		}
	}
	return area
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func maxAreaOptimized(height []int) int {
	max, start, end := 0, 0, len(height)-1
	for start < end {
		width := end - start
		high := 0
		if height[start] < height[end] {
			high = height[start]
			start++
		} else {
			high = height[end]
			end--
		}
		temp := width * high
		if temp > max {
			max = temp
		}
	}
	return max
}
