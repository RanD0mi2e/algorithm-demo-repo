package main

import (
	"fmt"
	"sort"
)

func main() {
	input := []int{0, 0, 0, 0}
	output := threeSum2(input)
	fmt.Printf("result: %v", output)
}

// 固定中间指针，左右指针组成双指针遍历
func threeSum(nums []int) [][]int {
	length := len(nums) - 1
	result := [][]int{}
	// 排序
	sort.Ints(nums)
	for i := 1; i < length; i++ {
		start, end := 0, length
		// i指针相同，跳过本次循环，避免添加重复数据
		if i > 1 && nums[i] == nums[i-1] {
			start = i - 1
		}
		// 双指针移动，保证左指针小于i指针，右指针大于i指针，避免重复数字被计算
		for start < i && end > i {
			// 左指针当前值和上一步值相同，直接跳过避免重复添加
			if start > 0 && nums[start] == nums[start-1] {
				start++
				continue
			}
			// 右指针同理
			if end < length && nums[end] == nums[end+1] {
				end--
				continue
			}
			// 计算三数之和，=0为结果项，>0说明和过大，右指针左移，<0说明和过小，左指针右移
			sum := nums[start] + nums[i] + nums[end]
			if sum == 0 {
				result = append(result, []int{nums[start], nums[i], nums[end]})
				start++
				end--
			} else if sum < 0 {
				start++
			} else {
				end--
			}
		}
	}
	return result
}

// 固定左指针，中间指针、右指针组成双指针遍历
func threeSum2(nums []int) [][]int {
	length := len(nums) - 1
	result := [][]int{}
	// 排序
	sort.Ints(nums)
	for i := 0; i < length-1; i++ {
		start, end := i+1, length
		// i指针相同，跳过本次循环，避免添加重复数据
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for start < end {
			if start > i+1 && nums[start] == nums[start-1] {
				start++
				continue
			}
			// 右指针同理
			if end < length && nums[end] == nums[end+1] {
				end--
				continue
			}
			// 计算三数之和，=0为结果项，>0说明和过大，右指针左移，<0说明和过小，左指针右移
			sum := nums[start] + nums[i] + nums[end]
			if sum == 0 {
				result = append(result, []int{nums[start], nums[i], nums[end]})
				start++
				end--
			} else if sum < 0 {
				start++
			} else {
				end--
			}
		}
	}
	return result
}
