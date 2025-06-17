package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func sortedArrayToBST(nums []int) *TreeNode {
	left := 0
	right := len(nums) - 1

	var recurve func(left, right int) *TreeNode
	recurve = func(left, right int) *TreeNode {
		if left > right {
			return nil
		}
		mid := (left + right) >> 1
		node := &TreeNode{
			Val:   nums[mid],
			Left:  recurve(left, mid-1),
			Right: recurve(mid+1, right),
		}
		return node
	}

	return recurve(left, right)
}

func main() {
	input := []int{-10, -3, 0, 5, 9}
	result := sortedArrayToBST(input)
	fmt.Printf("result: %T", result)
}
