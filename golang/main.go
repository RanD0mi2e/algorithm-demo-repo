package main

import "fmt"

func main() {

	var print []int = twoSum([]int{3, 2, 4}, 6)
	fmt.Println(print)
}

func twoSum(nums []int, target int) []int {
	var numsHashMap = make(map[int]int)
	for i, v := range nums {
		if p, exist := numsHashMap[target-v]; exist {
			return []int{p, i}
		}
		numsHashMap[v] = i
	}
	return nil
}
