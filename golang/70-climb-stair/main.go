package main

import "fmt"

func main() {
	input := 5
	output := climbStairs(input)
	fmt.Printf("result: %v", output)
}

// 爬楼梯-斐波那契数列变种
func climbStairs(n int) int {
	// dp[n] = dp[n-1] + dp[n-2]
	// dp[1] = 1, dp[2] = 2
	if n < 2 {
		return n
	}
	p1, p2 := 1, 2
	for i := 3; i <= n; i++ {
		current := p1 + p2
		p1 = p2
		p2 = current
	}

	return p2
}
