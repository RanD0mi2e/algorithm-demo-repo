package main

import "fmt"

type DictNode struct {
	children [26]*DictNode
	isEnd    bool
}

type WordDictionary struct {
	root *DictNode
}

func Constructor() WordDictionary {
	return WordDictionary{
		root: &DictNode{},
	}
}

func (this *WordDictionary) AddWord(word string) {
	node := this.root
	for _, char := range word {
		index := char - 'a'
		if node.children[index] == nil {
			node.children[index] = &DictNode{}
		}
		node = node.children[index]
	}
	node.isEnd = true
}

func (this *WordDictionary) Search(word string) bool {
	var recurve func(node *DictNode, word string) bool
	recurve = func(node *DictNode, word string) bool {
		if node == nil {
			return false
		}
		if len(word) == 0 {
			return node.isEnd
		}
		firstChar := word[0]
		if firstChar == '.' {
			for _, child := range node.children {
				if child != nil && recurve(child, word[1:]) {
					return true
				}
			}
			return false
		}

		index := firstChar - 'a'
		return node.children[index] != nil && recurve(node.children[index], word[1:])
	}
	return recurve(this.root, word)
}

func main() {
	dict := Constructor()
	dict.AddWord("app")
	res1 := dict.Search("apple")
	res2 := dict.Search("app")
	res3 := dict.Search("a.p")
	res4 := dict.Search(".pp")
	res5 := dict.Search("ap.")
	res6 := dict.Search("a..")
	fmt.Printf("res1: %t\n", res1)
	fmt.Printf("res2: %t\n", res2)
	fmt.Printf("res3: %t\n", res3)
	fmt.Printf("res4: %t\n", res4)
	fmt.Printf("res5: %t\n", res5)
	fmt.Printf("res6: %t\n", res6)
}
