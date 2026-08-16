// Произведение всех элементов кроме i-го, без деления.
//
// Два прохода: сначала записываем в result[i] произведение всех элементов
// слева от i (префикс), затем идём справа налево и домножаем на произведение
// всех элементов справа (суффикс). Время O(n), доп. память O(1)
// (не считая результата).
package main

import "fmt"

func productExceptSelf(nums []int) []int {
	result := make([]int, len(nums))

	prefix := 1
	for i, v := range nums {
		result[i] = prefix
		prefix *= v
	}

	suffix := 1
	for i := len(nums) - 1; i >= 0; i-- {
		result[i] *= suffix
		suffix *= nums[i]
	}

	return result
}

func main() {
	fmt.Println(productExceptSelf([]int{1, 2, 3}))    // [6 3 2]
	fmt.Println(productExceptSelf([]int{1, 2, 3, 4})) // [24 12 8 6]
}
