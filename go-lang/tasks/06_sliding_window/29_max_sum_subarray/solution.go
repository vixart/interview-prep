// Максимальная сумма подмассива длины k — окно фиксированного размера.
//
// Считаем сумму первого окна, дальше сдвигаем: прибавляем входящий элемент,
// вычитаем выходящий. Время O(n), память O(1).
package main

import "fmt"

func maxSumSubarray(nums []int, k int) int {
	if k <= 0 || len(nums) < k {
		return 0
	}

	sum := 0
	for i := 0; i < k; i++ {
		sum += nums[i]
	}

	best := sum
	for i := k; i < len(nums); i++ {
		sum += nums[i] - nums[i-k]
		if sum > best {
			best = sum
		}
	}

	return best
}

func main() {
	fmt.Println(maxSumSubarray([]int{2, 1, 5, 1, 3, 2}, 3))             // 9
	fmt.Println(maxSumSubarray([]int{2, 3, 4, 1, 5}, 2))                // 7
	fmt.Println(maxSumSubarray([]int{1, 4, 2, 10, 23, 3, 1, 0, 20}, 4)) // 39
}
