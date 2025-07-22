package main

func maxSubArr(nums []int) int {
	length := len(nums)
	currentSum, globalSum := nums[0], nums[0]
	for i := 1; i < length; i++ {
		currentSum = max(nums[i], currentSum+nums[i])
		globalSum = max(currentSum, globalSum)
	}

	return globalSum
}
