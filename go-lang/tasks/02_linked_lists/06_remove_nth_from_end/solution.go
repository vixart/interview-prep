// Удаление N-го элемента с конца за один проход.
//
// Два указателя: fast уходит вперёд на n узлов, затем fast и slow идут
// вместе. Когда fast достигает конца, slow стоит перед удаляемым узлом.
// Фиктивная голова (dummy) избавляет от частного случая удаления первого
// элемента. Время O(n), память O(1).
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}

	fast, slow := dummy, dummy
	for i := 0; i < n; i++ {
		fast = fast.Next
	}
	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}
	slow.Next = slow.Next.Next

	return dummy.Next
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
	fmt.Println(toSlice(removeNthFromEnd(fromSlice(1, 2, 3, 4, 5), 2))) // [1 2 3 5]
	fmt.Println(toSlice(removeNthFromEnd(fromSlice(1), 1)))             // []
	fmt.Println(toSlice(removeNthFromEnd(fromSlice(1, 2), 1)))          // [1]
}
