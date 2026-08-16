// Удаление всех нулей из слайса in-place с сохранением порядка.
//
// Классический приём "два указателя": k указывает, куда писать следующий
// ненулевой элемент. Время O(n), память O(1), порядок сохраняется.
package main

import "fmt"

func remove(in []int) []int {
	k := 0
	for _, v := range in {
		if v != 0 {
			in[k] = v
			k++
		}
	}

	return in[:k]
}

func main() {
	fmt.Println(remove([]int{}))              // []
	fmt.Println(remove([]int{0}))             // []
	fmt.Println(remove([]int{1, 0, 0, 2}))    // [1 2]
	fmt.Println(remove([]int{0, 0, 1, 2, 3})) // [1 2 3]
	fmt.Println(remove([]int{1, 2, 3}))       // [1 2 3]
}
