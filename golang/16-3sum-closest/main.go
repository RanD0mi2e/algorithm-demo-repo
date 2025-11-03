package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	input := []int{4, 0, 5, -5, 3, 3, 0, -4, -5}
	output := threeSumClosest(input, -2)
	fmt.Println(output)
}

func threeSumClosest(nums []int, target int) int {
	n, result, diff := len(nums), 0, math.MaxInt32
	sort.Ints(nums)
	if n <= 2 {
		return result
	}
	for i := 0; i < n-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j, k := i+1, n-1; j < k; {
			sum := nums[i] + nums[j] + nums[k]
			if abs(sum-target) < diff {
				result, diff = sum, abs(sum-target)
			}
			if sum == target {
				return result
			} else if sum < target {
				j++
			} else {
				k--
			}
		}
	}

	return result
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
