package main

import "fmt"

func main() {
	input := []int{-1, 2, 1}
	output := subarraySum(input, 2)
	fmt.Println(output)
}

func subarraySum(nums []int, k int) int {
	count := 0
	countSum := 0

	preSumMap := make(map[int]int)
	preSumMap[0] = 1
	for _, num := range nums {
		countSum += num
		requireSum := countSum - k

		if freq, found := preSumMap[requireSum]; found {
			count += freq
		}

		preSumMap[countSum]++
	}

	return count
}
