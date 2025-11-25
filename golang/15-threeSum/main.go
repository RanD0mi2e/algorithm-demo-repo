package main

import (
	"sort"
)

func threeSum(nums []int) [][]int {
	lens := len(nums)
	var result [][]int
	sort.Ints(nums)
	for i := 0; i < lens-2; i++ {
		if nums[i] > 0 {
			break
		}
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		j, k := i+1, lens-1
		for j < k {
			if nums[i]+nums[j]+nums[k] == 0 {
				result = append(result, []int{nums[i], nums[j], nums[k]})
				for j < k && nums[j] == nums[j+1] {
					j += 1
				}
				for j < k && nums[k] == nums[k-1] {
					k -= 1
				}
				j++
				k--
			} else if nums[i]+nums[j]+nums[k] < 0 {
				j++
			} else if nums[i]+nums[j]+nums[k] > 0 {
				k--
			}
		}
	}
	return result
}
