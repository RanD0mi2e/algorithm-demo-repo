package main

import (
	"fmt"
)

func main() {
	input := 25
	output := trailingZeroes(input)
	fmt.Printf("result: %v", output)
}

func trailingZeroes(n int) int {
	ans := 0
	for i := 5; n/i > 0; i *= 5 {
		ans += n / i
	}
	return ans
}
