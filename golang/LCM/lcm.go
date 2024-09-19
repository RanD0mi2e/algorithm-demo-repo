package LCM

import "algorithm/GCD"

// LCM(a, b) = (a * b) / GCD(a, b)
func LowestCommonMultiple(a, b int) int {
	return a * b / GCD.GreatestCommonDivisor(a, b)
}

func LowestCommonMultipleForNums(nums ...int) int {
	length := len(nums)
	result := nums[0]
	for i := 1; i < length; i++ {
		result = LowestCommonMultiple(result, nums[i])
	}
	return result
}
