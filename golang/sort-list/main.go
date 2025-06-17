package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	input := &ListNode{
		Val: 4,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 1,
				Next: &ListNode{
					Val:  3,
					Next: nil,
				},
			},
		},
	}
	result := sortList(input)
	fmt.Printf("result: %T", result)
}

func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	mid := findMiddle(head)
	rightHead := mid.Next
	mid.Next = nil

	leftLink := sortList(head)
	rightLink := sortList(rightHead)

	return merge(leftLink, rightLink)
}

func findMiddle(node *ListNode) *ListNode {
	if node == nil {
		return nil
	}
	slow, fast := node, node.Next
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

func merge(link1, link2 *ListNode) *ListNode {
	dummy := &ListNode{}
	current := dummy

	for link1 != nil && link2 != nil {
		if link1.Val > link2.Val {
			current.Next = link2
			link2 = link2.Next
		} else {
			current.Next = link1
			link1 = link1.Next
		}
		current = current.Next
	}

	if link1 != nil {
		current.Next = link1
	} else {
		current.Next = link2
	}

	return dummy.Next
}
