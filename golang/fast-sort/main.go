package main

import "fmt"

func main() {
	input := []int{4, 3, 2, 1}
	quickSort(input)
	fmt.Printf("result: %v", input)

}

func quickSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	stack := make([][2]int, 0)
	stack = append(stack, [2]int{0, len(arr) - 1})

	for len(stack) > 0 {
		n := len(stack) - 1
		low, high := stack[n][0], stack[n][1]
		stack = stack[:n]

		if low < high {
			p := partition(arr, low, high)

			if p-1 > low {
				stack = append(stack, [2]int{low, p - 1})
			}
			if p+1 < high {
				stack = append(stack, [2]int{p + 1, high})
			}
		}
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i] // 把 pivot 放到中间
	return i
}
