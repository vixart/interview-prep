// Контейнер с максимальной водой — два указателя с краёв.
//
// Площадь ограничена меньшей из двух линий, поэтому сдвигать имеет смысл
// именно её: сдвиг большей линии ширину уменьшит, а высоту не увеличит.
// Время O(n), память O(1).
package main

import "fmt"

func maxArea(height []int) int {
	best := 0
	l, r := 0, len(height)-1

	for l < r {
		h := min(height[l], height[r])
		if area := (r - l) * h; area > best {
			best = area
		}

		if height[l] < height[r] {
			l++
		} else {
			r--
		}
	}

	return best
}

func main() {
	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7})) // 49
	fmt.Println(maxArea([]int{1, 1}))                      // 1
	fmt.Println(maxArea([]int{4, 3, 2, 1, 4}))             // 16
	fmt.Println(maxArea([]int{1, 2, 1}))                   // 2
}
