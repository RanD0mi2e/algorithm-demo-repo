package main

import "fmt"

func main() {
	input1 := 121
	input2 := -121
	output1 := isPalindrome(input1)
	output2 := isPalindrome(input2)
	fmt.Printf("result1: %v\n", output1)
	fmt.Printf("result1: %v", output2)
}

func isPalindrome(x int) bool {
	// old
	// if x < 0 {
	// 	return false
	// }
	// var nums []int
	// for x > 0 {
	// 	nums = append(nums, x%10)
	// 	x = x / 10
	// }
	// left, right := 0, len(nums)-1
	// for left <= right {
	// 	if nums[left] != nums[right] {
	// 		return false
	// 	}
	// 	left++
	// 	right--
	// }
	// return true

	// new
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	reverseNum := 0
	for x > reverseNum {
		reverseNum = reverseNum*10 + x%10
		x /= 10
	}
	return x == reverseNum || x == reverseNum/10
}
