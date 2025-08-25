package main

import "fmt"

func main() {
	fmt.Printf("result: %v", myPow(2.00000, 7))
}

func myPow(x float64, n int) float64 {
	result := float64(1)
	// 处理负数
	if n < 0 {
		x = 1 / x
		n = -n
	}
	for n > 0 {
		if n&1 == 1 {
			result *= x
		}
		x *= x
		n >>= 1
	}
	return result
}
