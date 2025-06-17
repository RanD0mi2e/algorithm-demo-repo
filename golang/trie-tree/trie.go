package main

import "fmt"

type TrieNode struct {
	children [26]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{},
	}
}

func (t *Trie) insert(word string) {
	node := t.root
	for _, char := range word {
		index := char - 'a'
		if node.children[index] == nil {
			node.children[index] = &TrieNode{}
		}
		node = node.children[index]
	}
	node.isEnd = true
}

func (t *Trie) search(word string) bool {
	node := t.root
	for _, char := range word {
		index := char - 'a'
		if node.children[index] == nil {
			return false
		}
		node = node.children[index]
	}
	return node.isEnd
}

func (t *Trie) startsWith(word string) bool {
	node := t.root
	for _, char := range word {
		index := char - 'a'
		if node.children[index] == nil {
			return false
		}
		node = node.children[index]
	}
	return true
}

func main() {
	input := "apple"
	trie := NewTrie()
	trie.insert(input)
	res1 := trie.search("app")
	res2 := trie.startsWith("app")
	fmt.Printf("字典是否存在app: %t,\n字典是否存在以app开头的词: %t", res1, res2)
}
