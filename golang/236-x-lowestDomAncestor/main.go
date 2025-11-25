package main

import "fmt"

// mock dom
type DOMNode struct {
	ID       string
	Children []*DOMNode
}

func findLowestDOMAncestor(root, p, q *DOMNode) *DOMNode {
	if root == nil {
		return nil
	}

	if root == p || root == q {
		return root
	}

	var tempResult *DOMNode
	matchNodeNums := 0
	for _, node := range root.Children {
		result := findLowestDOMAncestor(node, p, q)
		if result != nil {
			tempResult = result
			matchNodeNums += 1
		}
	}

	if matchNodeNums >= 2 {
		return root
	}
	if matchNodeNums == 1 {
		return tempResult
	}
	return nil
}

func main() {
	p1 := &DOMNode{ID: "p1"}
	p2 := &DOMNode{ID: "p2"}
	div1 := &DOMNode{ID: "div1", Children: []*DOMNode{p1, p2}}

	span1 := &DOMNode{ID: "span1"}

	img1 := &DOMNode{ID: "img1"}
	div2 := &DOMNode{ID: "div2", Children: []*DOMNode{img1}}

	root := &DOMNode{ID: "body", Children: []*DOMNode{div1, span1, div2}}

	// 测试 1: 跨分支查询 (p1 和 img1) -> 应该是 body
	lca1 := findLowestDOMAncestor(root, p1, img1)
	fmt.Printf("LCA of p1 and img1: %s\n", lca1.ID)

	// 测试 2: 同分支查询 (p1 和 p2) -> 应该是 div1
	lca2 := findLowestDOMAncestor(root, p1, p2)
	fmt.Printf("LCA of p1 and p2: %s\n", lca2.ID)

	// 测试 3: 包含关系 (div1 和 p1) -> 应该是 div1
	lca3 := findLowestDOMAncestor(root, div1, p1)
	fmt.Printf("LCA of div1 and p1: %s\n", lca3.ID)
}
