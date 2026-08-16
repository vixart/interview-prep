// Удаление дубликатов из отсортированного массива in-place.
//
// Медленный указатель k показывает конец "чистой" части массива.
// Быстрый идёт по всем элементам и переносит значение, только если оно
// отличается от последнего записанного. Время O(n), память O(1).
package main

import "fmt"

func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	k := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[k-1] {
			nums[k] = nums[i]
			k++
		}
	}

	return k
}

func main() {
	a := []int{1, 1, 2}
	n := removeDuplicates(a)
	fmt.Println(n, a[:n]) // 2 [1 2]

	b := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	n = removeDuplicates(b)
	fmt.Println(n, b[:n]) // 5 [0 1 2 3 4]
}
