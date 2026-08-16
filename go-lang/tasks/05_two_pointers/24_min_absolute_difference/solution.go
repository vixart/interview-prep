// Наименьшая разница по модулю между элементами двух массивов.
//
// Сортируем оба массива и идём двумя указателями: сдвигаем тот, чей текущий
// элемент меньше — только так разница может уменьшиться.
// Время O(n log n + m log m), память O(1).
package main

import (
	"fmt"
	"sort"
)

func minAbsoluteDifference(arr1, arr2 []int) int {
	a := append([]int(nil), arr1...) // не портим входные слайсы
	b := append([]int(nil), arr2...)
	sort.Ints(a)
	sort.Ints(b)

	best := -1
	for i, j := 0, 0; i < len(a) && j < len(b); {
		d := a[i] - b[j]
		if d < 0 {
			d = -d
		}
		if best == -1 || d < best {
			best = d
		}

		if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}

	return best
}

func main() {
	fmt.Println(minAbsoluteDifference(
		[]int{1, 10, 15, 4, 20},
		[]int{3, 16, 5, 7},
	)) // 1 (4 и 5, или 15 и 16)
}
