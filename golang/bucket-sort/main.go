package main

import (
	"fmt"
	"sort"
)

func main() {
	input := []int{1, 7, 3, 5, 2, 9}
	output := bucketSort(input)
	fmt.Printf("result: %v", output)

	input2 := []float64{1.0, 7.1, 3.2, 5.3, 2.4, 9.5}
	output2 := bucketSortFloat(input2)
	fmt.Printf("result2: %v", output2)
}

// int type sort
func bucketSort(arr []int) []int {
	min, max := arr[0], arr[0]
	// find min, max num in arr
	for _, value := range arr {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	result := []int{}
	bucketCount := len(arr)
	if bucketCount == 0 {
		return result
	}

	// create buckets
	buckets := make([][]int, bucketCount)
	for i := range buckets {
		buckets[i] = make([]int, 0)
	}

	for _, value := range arr {
		bucketIndex := 0
		if max-min > 0 {
			bucketIndex = (value - min) * (bucketCount - 1) / (max - min)
		}
		buckets[bucketIndex] = append(buckets[bucketIndex], value)
	}

	for i := range buckets {
		if len(buckets[i]) > 0 {
			sort.Ints(buckets[i])
		}
	}

	for _, bucket := range buckets {
		result = append(result, bucket...)
	}

	return result
}

// flot type sort
func bucketSortFloat(arr []float64) []float64 {
	if len(arr) <= 1 {
		return arr
	}

	result := []float64{}

	min, max := arr[0], arr[0]
	for _, value := range arr {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	bucketCount := len(arr)
	buckets := make([][]float64, bucketCount)
	for i := range buckets {
		buckets[i] = make([]float64, 0)
	}

	for _, value := range arr {
		bucketIndex := 0
		if max-min > 0 {
			bucketIndex = int((value - min) / (max - min) * float64((bucketCount - 1)))
		}
		buckets[bucketIndex] = append(buckets[bucketIndex], value)
	}

	for _, bucket := range buckets {
		result = append(result, bucket...)
	}

	return result
}
