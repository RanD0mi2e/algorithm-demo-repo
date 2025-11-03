package main

import "math"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func main() {
	root := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
		},
		Right: &TreeNode{
			Val: 3,
		},
	}
	res := maxPathSum(root)
	println(res)
}

func maxPathSum(root *TreeNode) int {
	maxSum := math.MinInt
	var postOrderTraversal func(*TreeNode) int
	postOrderTraversal = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		leftVal := max(postOrderTraversal(root.Left), 0)
		rightVal := max(postOrderTraversal(root.Right), 0)
		candidate := root.Val + leftVal + rightVal
		maxSum = max(maxSum, candidate)
		return root.Val + max(leftVal, rightVal)
	}
	postOrderTraversal(root)
	return maxSum
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
