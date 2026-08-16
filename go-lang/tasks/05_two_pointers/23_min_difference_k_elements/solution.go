// Минимальная разница между максимумом и минимумом из k выбранных элементов.
//
// После сортировки оптимальные k элементов всегда стоят подряд, поэтому
// достаточно пройтись окном размера k и взять минимум nums[i+k-1]-nums[i].
// Время O(n log n), память O(1).
package main

import (
	"fmt"
	"sort"
)

func minDifference(nums []int, k int) int {
	if k <= 1 || len(nums) < k {
		return 0
	}

	sort.Ints(nums)

	best := nums[k-1] - nums[0]
	for i := k; i < len(nums); i++ {
		if d := nums[i] - nums[i-k+1]; d < best {
			best = d
		}
	}

	return best
}

func main() {
	fmt.Println(minDifference([]int{1, 3, 5, 7, 9}, 3))  // 4
	fmt.Println(minDifference([]int{5, 1, 100, 102}, 2)) // 2 (100 и 102)
}
