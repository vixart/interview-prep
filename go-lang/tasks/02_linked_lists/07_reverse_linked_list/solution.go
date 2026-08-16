// Разворот односвязного списка: итеративно и рекурсивно.
//
// Итеративный вариант — O(n) по времени, O(1) по памяти:
// перекидываем указатель Next каждого узла на предыдущий.
// Рекурсивный — O(n) по времени, O(n) по памяти (стек вызовов).
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode

	for head != nil {
		next := head.Next
		head.Next = prev
		prev = head
		head = next
	}

	return prev
}

func reverseListRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := reverseListRecursive(head.Next)
	head.Next.Next = head
	head.Next = nil

	return newHead
}

func fromSlice(vals ...int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func toSlice(head *ListNode) []int {
	var out []int
	for ; head != nil; head = head.Next {
		out = append(out, head.Val)
	}
	return out
}

func main() {
	fmt.Println(toSlice(reverseList(fromSlice(1, 2, 3, 4, 5)))) // [5 4 3 2 1]
	fmt.Println(toSlice(reverseListRecursive(fromSlice(1, 2)))) // [2 1]
	fmt.Println(toSlice(reverseList(nil)))                      // []
}
