// K ближайших элементов к arr[index] в отсортированном массиве.
//
// Окно [l, r] начинается с самого элемента и расширяется k-1 раз:
// на каждом шаге берём более близкого соседа; при равных расстояниях —
// правого (он больше, как требует условие). Время O(k), память O(1).
package main

import "fmt"

func findKClosest(arr []int, index int, k int) []int {
	if k == 0 {
		return []int{}
	}

	l, r := index, index
	for r-l+1 < k {
		switch {
		case l == 0:
			r++
		case r == len(arr)-1:
			l--
		case arr[index]-arr[l-1] < arr[r+1]-arr[index]:
			l--
		default: // правый ближе либо расстояния равны -> берём больший
			r++
		}
	}

	return arr[l : r+1]
}

func main() {
	fmt.Println(findKClosest([]int{2, 3, 5, 7, 11}, 3, 2))    // [5 7]
	fmt.Println(findKClosest([]int{4, 12, 15, 15, 24}, 1, 3)) // [12 15 15]
}
