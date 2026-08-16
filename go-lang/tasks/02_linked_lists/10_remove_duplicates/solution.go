// Удаление дубликатов из отсортированного списка.
//
// Список отсортирован, поэтому дубликаты стоят подряд: если значение
// следующего узла равно текущему — выкидываем следующий узел.
// Время O(n), память O(1).
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	for cur := head; cur != nil && cur.Next != nil; {
		if cur.Val == cur.Next.Val {
			cur.Next = cur.Next.Next
		} else {
			cur = cur.Next
		}
	}

	return head
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
	fmt.Println(toSlice(deleteDuplicates(fromSlice(1, 1, 2))))       // [1 2]
	fmt.Println(toSlice(deleteDuplicates(fromSlice(1, 1, 2, 3, 3)))) // [1 2 3]
	fmt.Println(toSlice(deleteDuplicates(fromSlice(1, 1, 1, 1))))    // [1]
}
