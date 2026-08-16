// Вращение массива вправо на k позиций через три разворота.
//
// [1,2,3,4,5,6,7], k=3:
//  1. развернуть всё:        [7,6,5,4,3,2,1]
//  2. развернуть первые k:   [5,6,7,4,3,2,1]
//  3. развернуть остаток:    [5,6,7,1,2,3,4]
//
// Время O(n), память O(1), in-place.
package main

import "fmt"

func rotate(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n

	reverse(nums)
	reverse(nums[:k])
	reverse(nums[k:])
}

func reverse(s []int) {
	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		s[l], s[r] = s[r], s[l]
	}
}

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	rotate(a, 3)
	fmt.Println(a) // [5 6 7 1 2 3 4]

	b := []int{-1, -100, 3, 99}
	rotate(b, 2)
	fmt.Println(b) // [3 99 -1 -100]
}
