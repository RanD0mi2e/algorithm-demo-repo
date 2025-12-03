package main

import "fmt"

// 41. 缺失的第一个正数

// 给你一个未排序的整数数组 nums ，请你找出其中没有出现的最小的正整数。

// 请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。

// 示例 1：

// 输入：nums = [1,2,0]
// 输出：3
// 解释：范围 [1,2] 中的数字都在数组中。
// 示例 2：

// 输入：nums = [3,4,2,1]
// 输出：2
// 解释：1 在数组中，但 2 没有。
// 示例 3：

// 输入：nums = [7,8,9,11,12]
// 输出：1
// 解释：最小的正数 1 没有出现。

func main() {
	input := []int{3, 2, 3}
	output := firstMissingPositive(input)
	fmt.Println(output)
}

func firstMissingPositive(nums []int) int {
	n := len(nums)
	for i := range n {
		for nums[i] > 0 && nums[i] < n && nums[nums[i]-1] != nums[i] {
			targetIndex := nums[i] - 1
			nums[i], nums[targetIndex] = nums[targetIndex], nums[i]
		}
	}

	for i := range n {
		if nums[i]-1 != i {
			return i + 1
		}
	}

	return n + 1
}
