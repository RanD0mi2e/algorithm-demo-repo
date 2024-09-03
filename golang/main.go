package main

import (
	AVLTree "algorithm/aVL-Tree"
	"encoding/json"
	"fmt"
)

func main() {

	tree := AVLTree.NewAVLTree()
	tree.Insert(1)
	tree.Insert(3)
	tree.Insert(5)
	tree.Insert(2)
	tree.Delete(1)
	bytes, err := json.Marshal(tree)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
