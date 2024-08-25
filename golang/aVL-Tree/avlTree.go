package avlTree

type AVLTree struct {
	root *TreeNode
}

type TreeNode struct {
	Val    int
	Left   *TreeNode
	Right  *TreeNode
	Height int
}

func NewAVLTree() *AVLTree {
	return &AVLTree{
		root: nil,
	}
}

func NewTreeNode(value int) *TreeNode {
	return &TreeNode{
		Val:    value,
		Left:   nil,
		Right:  nil,
		Height: 0,
	}
}

func (t *AVLTree) height(node *TreeNode) int {
	if node != nil {
		return node.Height
	}

	return -1
}

func (t *AVLTree) updateHeight(node *TreeNode) {
	lh := t.height(node.Left)
	rh := t.height(node.Right)

	if lh > rh {
		node.Height = lh + 1
	} else {
		node.Height = rh + 1
	}
}

func (t *AVLTree) balanceFactor(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return t.height(node.Left) - t.height(node.Right)
}

func (t *AVLTree) rightRotate(node *TreeNode) *TreeNode {
	child := node.Left
	grandChild := child.Right
	child.Right = node
	node.Left = grandChild
	t.updateHeight(node)
	t.updateHeight(child)
	return child
}

func (t *AVLTree) leftRotate(node *TreeNode) *TreeNode {
	child := node.Right
	grandChild := child.Left
	child.Left = node
	node.Right = grandChild
	t.updateHeight(node)
	t.updateHeight(child)
	return child
}

func (t *AVLTree) rotate(node *TreeNode) *TreeNode {
	bf := t.balanceFactor(node)
	if bf > 1 {
		if t.balanceFactor(node.Left) > 0 {
			return t.rightRotate(node)
		} else {
			node.Left = t.leftRotate(node.Left)
			return t.rightRotate(node)
		}
	}

	if bf < -1 {
		if t.balanceFactor(node.Right) < 0 {
			return t.leftRotate(node)
		} else {
			node.Right = t.rightRotate(node.Right)
			return t.leftRotate(node)
		}
	}

	return node
}
