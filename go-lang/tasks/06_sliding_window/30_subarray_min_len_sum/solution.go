// Минимальная длина подмассива с суммой >= k — динамическое окно.
//
// Правая граница расширяет окно, увеличивая сумму; как только sum >= k,
// сжимаем окно слева, фиксируя минимум длины. Работает потому, что все
// числа положительные (сумма монотонна). Время O(n), память O(1).
package main

import "fmt"

func minSubArrayLen(nums []int, k int) int {
	best := 0
	sum, l := 0, 0

	for r, v := range nums {
		sum += v

		for sum >= k {
			if best == 0 || r-l+1 < best {
				best = r - l + 1
			}
			sum -= nums[l]
			l++
		}
	}

	return best
}

func main() {
	fmt.Println(minSubArrayLen([]int{2, 3, 1, 2, 4, 3}, 7))        // 2 ([4,3])
	fmt.Println(minSubArrayLen([]int{1, 4, 4}, 4))                 // 1 ([4])
	fmt.Println(minSubArrayLen([]int{1, 1, 1, 1, 1, 1, 1, 1}, 11)) // 0
}
