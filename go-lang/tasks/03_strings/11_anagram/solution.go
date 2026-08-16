// Анаграммы: сравнение частот символов через хеш-таблицу.
//
// Считаем частоты рун первой строки, вычитаем частоты второй.
// Работает с любым Unicode (руны, а не байты), регистр и пробелы значимы.
// Время O(n), память O(k) — k уникальных символов.
//
// Альтернатива — отсортировать руны обеих строк и сравнить: O(n log n).
package main

import "fmt"

func IsAnagram(a, b string) bool {
	counts := make(map[rune]int)

	for _, r := range a {
		counts[r]++
	}
	for _, r := range b {
		counts[r]--
		if counts[r] < 0 {
			return false
		}
	}
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println(IsAnagram("лапоть", "пальто")) // true
	fmt.Println(IsAnagram("listen", "silent")) // true
	fmt.Println(IsAnagram("hello", "world"))   // false
	fmt.Println(IsAnagram("", ""))             // true
	fmt.Println(IsAnagram("a", "b"))           // false
}
