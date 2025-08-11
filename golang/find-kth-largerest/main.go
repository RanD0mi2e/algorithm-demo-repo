package main

import (
	"fmt"
	"math/rand"
)

func main() {
	result := findKthLargest2([]int{3, 2, 1, 5, 6, 4}, 2)
	fmt.Printf("result: %v", result)
}

func findKthLargest(nums []int, k int) int {
	result := []int{}

	up := func(index int) {
		for index > 0 {
			parent := (index - 1) >> 1
			if result[parent] < result[index] {
				result[parent], result[index] = result[index], result[parent]
			}
			index = parent
		}
	}

	down := func(index int) {
		n := len(result)
		for {
			left := index*2 + 1
			right := index*2 + 2
			if left >= n {
				break
			}
			maxChild := left
			if right < n && result[maxChild] < result[right] {
				maxChild = right
			}
			if result[maxChild] < result[index] {
				break
			}
			result[index], result[maxChild] = result[maxChild], result[index]
			index = maxChild
		}
	}

	insert := func(val int) {
		result = append(result, val)
		idx := len(result) - 1
		up(idx)
	}

	heapify := func(arr []int) {
		for _, val := range arr {
			insert(val)
		}
	}

	heapify(nums)
	for i := 0; i < k-1; i++ {
		result[0] = result[len(result)-1]
		result = result[:len(result)-1]
		down(0)
	}

	return result[0]
}

func findKthLargest2(nums []int, k int) int {
	partion := func(nums []int, left, right int, pivotIndex int) int {
		pivotValue := nums[pivotIndex]
		nums[right], nums[pivotIndex] = nums[pivotIndex], nums[right]
		storeIndex := left
		for i := left; i < right; i++ {
			if nums[i] < pivotValue {
				nums[i], nums[storeIndex] = nums[storeIndex], nums[i]
				storeIndex++
			}
		}
		// 将基准元素放到正确位置
		nums[right], nums[storeIndex] = nums[storeIndex], nums[right]
		return storeIndex
	}

	var quickSelect func(nums []int, left, right int, targetIndex int) int
	quickSelect = func(nums []int, left, right int, targetIndex int) int {
		if left == right {
			return nums[left]
		}
		pivotIndex := left + rand.Intn(right-left+1)
		pivotIndex = partion(nums, left, right, pivotIndex)
		if pivotIndex == targetIndex {
			return nums[targetIndex]
		} else if pivotIndex > targetIndex {
			return quickSelect(nums, left, pivotIndex-1, targetIndex)
		} else {
			return quickSelect(nums, pivotIndex+1, right, targetIndex)
		}
	}

	targetIndex := len(nums) - k
	return quickSelect(nums, 0, len(nums)-1, targetIndex)
}
