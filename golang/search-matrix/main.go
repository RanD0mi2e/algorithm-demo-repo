package main

import "fmt"

func main() {
	input := [][]int{{1}}
	// output := searchMatrix(input, 0)
	output := searchMatrix2(input, 0)
	fmt.Printf("result: %t", output)
}

// 双层二分查找
func searchMatrix(matrix [][]int, target int) bool {
	n := len(matrix)
	if n == 0 {
		return false
	}

	left, right := 0, n-1
	for left <= right {
		outMid := (left + right) >> 1
		length := len(matrix[outMid])
		inLeft, inRight := 0, length-1
		for inLeft <= inRight {
			if matrix[outMid][inLeft] > target {
				right = outMid - 1
				break
			}
			if matrix[outMid][inRight] < target {
				left = outMid + 1
				break
			}
			inMid := (inLeft + inRight) >> 1
			if matrix[outMid][inMid] > target {
				inRight = inMid - 1
			} else if matrix[outMid][inMid] < target {
				inLeft = inMid + 1
			} else {
				return true
			}
		}
	}

	return false
}

// 把整个矩阵当作长数组单次二分查找
func searchMatrix2(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	left, right := 0, m*n-1
	for left <= right {
		mid := (left + right) >> 1
		if target > matrix[mid/n][mid%n] {
			left = mid + 1
		} else if target < matrix[mid/n][mid%n] {
			right = mid - 1
		} else {
			return true
		}
	}
	return false
}
