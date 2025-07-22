package main

import (
	"fmt"
)

func solveNQueens(n int) [][]string {
	// 1.创建棋盘
	result := [][]string{}
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}
	// 2.从第一行开始放置皇后
	backtrack(board, 0, &result)
	return result
}

func backtrack(board [][]byte, row int, result *[][]string) {
	// 放置完n个皇后，递归终点
	if row == len(board) {
		*result = append(*result, constructor(board))
		return
	}
	// 3.当前行每个位置都尝试放置皇后
	for col := 0; col < len(board); col++ {
		// 4.当前行、列所在位置是否能放置新的皇后
		if isValid(board, row, col) {
			board[row][col] = 'Q'
			// 5.寻找下一个皇后放置位置
			backtrack(board, row+1, result)
			// 6.回溯
			board[row][col] = '.'
		}
	}
}

func isValid(board [][]byte, row, col int) bool {
	n := len(board)
	// 同一列不允许多个皇后
	for i := range row {
		if board[i][col] == 'Q' {
			return false
		}
	}
	// 左上对角线不允许多个皇后
	for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if board[i][j] == 'Q' {
			return false
		}
	}
	// 右上对角线不允许多个皇后
	for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
		if board[i][j] == 'Q' {
			return false
		}
	}
	return true
}

// 构造返回结果
func constructor(board [][]byte) []string {
	result := make([]string, len(board))
	for i := range board {
		result[i] = string(board[i])
	}
	return result
}

func main() {
	n := 5 // 可以修改为任意N值
	solutions := solveNQueens(n)
	fmt.Printf("Total %d solutions for %d-Queens problem:\n", len(solutions), n)
	for i, solution := range solutions {
		fmt.Printf("Solution %d:\n", i+1)
		for _, row := range solution {
			fmt.Println(row)
		}
		fmt.Println()
	}
}
