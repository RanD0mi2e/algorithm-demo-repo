package main

import (
	"fmt"
	"sort"
)

func main() {
	input1 := []int{1, 2, 3, 1} // 1 1 2 3
	input2 := []int{2, 2, 3}    // 2 2 3
	result := intersection(input1, input2)
	fmt.Println("result: ", result)
}

func intersection(nums1 []int, nums2 []int) []int {
	var result []int

	sort.Ints(nums1)
	sort.Ints(nums2)

	i, j := 0, 0
	for i < len(nums1) && j < len(nums2) {
		if nums1[i] < nums2[j] {
			i++
		} else if nums1[i] > nums2[j] {
			j++
		} else {
			if len(result) == 0 || result[len(result)-1] != nums1[i] {
				result = append(result, nums1[i])
			}
			i++
			j++
		}
	}
	return result
}

func max(i, j int) int {
	if i < j {
		return j
	}
	return i
}
