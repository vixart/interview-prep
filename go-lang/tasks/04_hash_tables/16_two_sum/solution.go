// Two Sum: есть ли в массиве пара с суммой target.
//
// Один проход с множеством "виденных" чисел: для текущего v проверяем,
// встречали ли мы target-v. Время O(n), память O(n) — вместо
// наивного двойного цикла O(n²).
package main

func HasSum(nums []int, target int) bool {
	seen := make(map[int]struct{}, len(nums))

	for _, v := range nums {
		if _, ok := seen[target-v]; ok {
			return true
		}
		seen[v] = struct{}{}
	}

	return false
}

func main() {
	println(HasSum([]int{10, 15, 3, 7}, 17))     // true (10 + 7 = 17)
	println(HasSum([]int{1, 2, 3, 4, 5, 6}, 12)) // false
	println(HasSum([]int{5, 5}, 10))             // true
}
