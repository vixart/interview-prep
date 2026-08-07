// Методы можно вызывать на NIL-приемнике: дерево строится от `var it *IntTree` (nil),
// Insert и Contains корректно обрабатывают it == nil. Это работает потому,
// что приемник указательный — вызов метода на nil-указателе не паникует сам по себе.
package main

import (
	"fmt"
)

type IntTree struct {
	val         int
	left, right *IntTree
}

func (it *IntTree) Insert(val int) *IntTree {
	if it == nil {
		// вызов метода на nil-указателе легален — приемник просто равен nil
		return &IntTree{val: val}
	}
	if val < it.val {
		it.left = it.left.Insert(val)
	} else if val > it.val {
		it.right = it.right.Insert(val)
	}
	return it
}

func (it *IntTree) Contains(val int) bool {
	switch {
	case it == nil:
		return false
	case val < it.val:
		return it.left.Contains(val)
	case val > it.val:
		return it.right.Contains(val)
	default:
		return true
	}
}

func main() {
	var it *IntTree
	// начинаем с nil-дерева...
	it = it.Insert(5)
	// ...и спокойно вызываем на нем метод
	it = it.Insert(3)
	it = it.Insert(10)
	it = it.Insert(2)
	fmt.Println(it.Contains(2))
	fmt.Println(it.Contains(12))
}
