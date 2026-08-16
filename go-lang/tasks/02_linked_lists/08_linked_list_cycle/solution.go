// Поиск цикла в списке — алгоритм Флойда («черепаха и заяц»).
//
// Медленный указатель шагает по одному узлу, быстрый — по два.
// Если цикл есть, быстрый обязательно догонит медленного внутри цикла;
// если нет — быстрый упрётся в nil. Время O(n), память O(1).
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func hasCycle(head *ListNode) bool {
	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}

	return false
}

// makeList строит список из vals; pos — индекс узла, на который замыкается
// хвост (-1 — без цикла).
func makeList(vals []int, pos int) *ListNode {
	dummy := &ListNode{}
	cur := dummy

	var nodes []*ListNode
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
		nodes = append(nodes, cur)
	}
	if pos >= 0 {
		cur.Next = nodes[pos]
	}

	return dummy.Next
}

func main() {
	fmt.Println(hasCycle(makeList([]int{3, 2, 0, -4}, 1))) // true
	fmt.Println(hasCycle(makeList([]int{1, 2}, 0)))        // true
	fmt.Println(hasCycle(makeList([]int{1}, -1)))          // false
}
