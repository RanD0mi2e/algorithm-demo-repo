package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root == p || root == q {
		return root
	}

	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)

	if left != nil && right != nil {
		return root
	}
	if left != nil {
		return left
	}
	if right != nil {
		return right
	}
	return nil
}

func main() {
	node7 := &TreeNode{Val: 7}
	node4 := &TreeNode{Val: 4}
	node6 := &TreeNode{Val: 6}
	node2 := &TreeNode{Val: 2, Left: node7, Right: node4}
	node0 := &TreeNode{Val: 0}
	node8 := &TreeNode{Val: 8}
	node5 := &TreeNode{Val: 5, Left: node6, Right: node2}
	node1 := &TreeNode{Val: 1, Left: node0, Right: node8}
	root := &TreeNode{Val: 3, Left: node5, Right: node1}
	fmt.Printf("LCA(5, 1) = %d\n", lowestCommonAncestor(root, node5, node1).Val)
}
