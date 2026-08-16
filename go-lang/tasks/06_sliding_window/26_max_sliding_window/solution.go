// Максимум в скользящем окне — монотонная убывающая deque.
//
// В deque храним ИНДЕКСЫ элементов, значения которых убывают.
// Голова deque — всегда максимум текущего окна.
//   - с головы выбрасываем индексы, вышедшие из окна;
//   - с хвоста выбрасываем все элементы, меньшие нового (они уже никогда
//     не станут максимумом).
//
// Каждый индекс входит и выходит из deque по одному разу -> O(n).
package main

import "fmt"

func maxSlidingWindow(nums []int, k int) []int {
	if k == 0 || len(nums) == 0 {
		return nil
	}

	result := make([]int, 0, len(nums)-k+1)
	deque := make([]int, 0, k) // индексы, nums по ним убывают

	for i, v := range nums {
		// убрать вышедший из окна индекс
		if len(deque) > 0 && deque[0] <= i-k {
			deque = deque[1:]
		}
		// убрать с хвоста всё, что меньше текущего
		for len(deque) > 0 && nums[deque[len(deque)-1]] < v {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)

		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

func main() {
	fmt.Println(maxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)) // [3 3 5 5 6 7]
	fmt.Println(maxSlidingWindow([]int{1}, 1))                        // [1]
}
