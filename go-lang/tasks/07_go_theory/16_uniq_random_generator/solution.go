// uniqN: слайс из n уникальных случайных чисел.
//
// map[int]struct{} используется как множество (пустая структура занимает
// 0 байт). Память под мапу и результат выделяем заранее.
// Помним: порядок обхода мапы в Go случайный.
package main

import (
	"fmt"
	"math/rand"
)

func uniqN(n int) []int {
	m := make(map[int]struct{}, n)

	for len(m) < n {
		m[rand.Int()] = struct{}{}
	}

	result := make([]int, 0, n)
	for k := range m {
		result = append(result, k)
	}

	return result
}

func main() {
	fmt.Println(uniqN(10))
	fmt.Println(len(uniqN(1000))) // 1000 — все уникальны
}
