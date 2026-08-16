// Слияние двух отсортированных связных списков без создания новых узлов.
//
// Фиктивная голова + хвостовой указатель: на каждом шаге пришиваем к хвосту
// меньший из двух текущих узлов. Остаток непустого списка пришиваем целиком.
// Время O(n+m), память O(1).
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			tail.Next = list1
			list1 = list1.Next
		} else {
			tail.Next = list2
			list2 = list2.Next
		}
		tail = tail.Next
	}

	if list1 != nil {
		tail.Next = list1
	} else {
		tail.Next = list2
	}

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
	fmt.Println(toSlice(mergeTwoLists(fromSlice(1, 2, 4), fromSlice(1, 3, 4)))) // [1 1 2 3 4 4]
	fmt.Println(toSlice(mergeTwoLists(nil, nil)))                               // []
	fmt.Println(toSlice(mergeTwoLists(nil, fromSlice(0))))                      // [0]
}
