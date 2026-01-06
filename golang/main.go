package main

import "fmt"

func main() {
	input := [][]int{
		{-1, 3},
	}
	// input := [][]int{
	// 	{1, 4, 7, 11, 15},
	// 	{2, 5, 8, 12, 19},
	// 	{3, 6, 9, 16, 22},
	// 	{10, 13, 14, 17, 24},
	// 	{18, 21, 23, 26, 30},
	// }
	output := searchMatrix(input, 3)
	fmt.Println("output: ", output)
}

func searchMatrix(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	row, col := 0, m-1
	for row < n && col >= 0 {
		if target > matrix[col][row] {
			row++
		} else if target < matrix[col][row] {
			col--
		} else {
			return true
		}
	}
	return false
}
