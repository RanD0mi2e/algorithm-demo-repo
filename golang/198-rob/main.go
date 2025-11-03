package main

import "fmt"

func main() {
	input := []int{0}
	output := rob(input)
	fmt.Printf("result: %v", output)
}

func rob(nums []int) int {
	// n := len(nums)
	// if n < 2 {
	// 	return nums[0]
	// }
	// dp := make([]int, len(nums))
	// dp[0] = nums[0]
	// dp[1] = max(nums[0], nums[1])
	// for i := 2; i < n; i++ {
	// 	dp[i] = max(dp[i-1], dp[i-2]+nums[i])
	// }
	// return dp[n-1]

	// 空间优化，当前状态只依赖前两个状态
	n := len(nums)
	if n < 2 {
		return nums[0]
	}
	p1, p2 := nums[0], max(nums[0], nums[1])
	for i := 2; i < n; i++ {
		current := max(p1+nums[i], p2)
		p1 = p2
		p2 = current
	}
	return p2
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
