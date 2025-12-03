package main

func main() {}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func leafSimilar(root1 *TreeNode, root2 *TreeNode) bool {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func() {
		dfs(root1, chan1)
		close(chan1)
	}()

	go func() {
		dfs(root2, chan2)
		close(chan2)
	}()

	for {
		v1, ok1 := <-chan1
		v2, ok2 := <-chan2

		if ok1 != ok2 {
			return false
		}

		if !ok1 {
			break
		}

		if v1 != v2 {
			return false
		}
	}

	return true
}

func dfs(node *TreeNode, ch chan int) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		ch <- node.Val
		return
	}

	dfs(node.Left, ch)
	dfs(node.Right, ch)
}
