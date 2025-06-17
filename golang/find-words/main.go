package main

import (
	"fmt"
	"strings"
)

type Trie struct {
	Children [26]*Trie
	Word     string
}

func NewTrie() *Trie {
	return &Trie{
		Children: [26]*Trie{},
		Word:     "",
	}
}

func (t *Trie) Insert(word string) {
	node := t
	for _, char := range word {
		index := char - 'a'
		if node.Children[index] == nil {
			node.Children[index] = NewTrie()
		}
		node = node.Children[index]
	}
	node.Word = word
}

func findWords(board [][]byte, words []string) []string {
	// 构建 trie 字典树
	trie := NewTrie()
	for _, word := range words {
		trie.Insert(word)
	}

	var result []string
	// 检查图的每一个位置是否是单词起点
	for col := range board {
		for row := range board[0] {
			bfs(board, col, row, trie, &result)
		}
	}

	return result
}

func bfs(board [][]byte, i, j int, node *Trie, result *[]string) {
	// 越界
	if i < 0 || i >= len(board) || j < 0 || j >= len(board[0]) {
		return
	}
	// 元素已访问过 或 字典中没有符合的字符
	char := board[i][j]
	index := char - 'a'
	if char == '#' || node.Children[index] == nil {
		return
	}
	// 添加匹配的单词后清空字典树中节点单词项，避免bfs重复遍历
	node = node.Children[index]
	if node.Word != "" {
		*result = append(*result, node.Word)
		node.Word = ""
	}
	// 已访问过
	board[i][j] = '#'
	direction := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, dir := range direction {
		bfs(board, i+dir[0], j+dir[1], node, result)
	}
	// 访问后回溯
	board[i][j] = char
}

func main() {
	board := [][]byte{
		{'a', 'o'},
	}
	words := []string{"oa", "ao", "o"}
	res := findWords(board, words)
	fmt.Printf("%s\n", strings.Join(res, ", "))
}
