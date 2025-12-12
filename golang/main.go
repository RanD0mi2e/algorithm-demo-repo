package main

import "fmt"

func main() {
	input := []int{0, -1}
	output := longestConsecutive(input)
	fmt.Println(output)
}

func longestConsecutive(nums []int) int {
	hm := make(map[int]bool)
	for _, num := range nums {
		hm[num] = true
	}

	maxLength := 0
	for num := range hm {
		if !hm[num-1] {
			length := 0
			for current := num; hm[current]; current++ {
				length += 1
			}
			maxLength = max(length, maxLength)
		}
	}

	return maxLength
}

func max(i, j int) int {
	if i < j {
		return j
	}
	return i
}
