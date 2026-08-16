// Find All Numbers Disappeared in an Array.
//
// Используем сам массив как хеш-таблицу: для каждого значения v помечаем
// ячейку с индексом v-1, делая её отрицательной. Индексы оставшихся
// положительных ячеек +1 — это пропавшие числа.
// Время O(n), доп. память O(1) (результат не считается).
package main

import "fmt"

func findDisappearedNumbers(nums []int) []int {
	for _, v := range nums {
		idx := abs(v) - 1
		if nums[idx] > 0 {
			nums[idx] = -nums[idx]
		}
	}

	var result []int
	for i, v := range nums {
		if v > 0 {
			result = append(result, i+1)
		}
	}

	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println(findDisappearedNumbers([]int{4, 3, 2, 7, 8, 2, 3, 1})) // [5 6]
	fmt.Println(findDisappearedNumbers([]int{1, 1}))                   // [2]
}
