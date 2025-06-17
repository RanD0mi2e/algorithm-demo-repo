package main

import "fmt"

func main() {
	input := [][]byte{{'A', 'B', 'C', 'E'}, {'S', 'F', 'C', 'S'}, {'A', 'D', 'E', 'E'}}
	result := exist(input, "ABCCED")
	fmt.Printf("result: %t", result)
}

func exist(board [][]byte, word string) bool {
	if len(board) == 0 || len(board[0]) == 0 || len(word) == 0 {
		return false
	}

	var bfs func(x, y int, index int) bool
	bfs = func(x, y int, index int) bool {
		if index == len(word) {
			return true
		}
		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) {
			return false
		}
		if board[x][y] != word[index] {
			return false
		}

		temp := board[x][y]
		board[x][y] = '#'

		direction := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, item := range direction {
			if bfs(x+item[0], y+item[1], index+1) {
				board[x][y] = temp
				return true
			}
		}

		board[x][y] = temp
		return false
	}

	for x, col := range board {
		for y := range col {
			if bfs(x, y, 0) {
				return true
			}
		}
	}

	return false
}
