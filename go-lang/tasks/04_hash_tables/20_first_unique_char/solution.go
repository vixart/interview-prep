// Первый уникальный символ — частотная таблица за два прохода.
//
// Первый проход считает частоты, второй находит первый символ с частотой 1.
// Время O(n), память O(1) (алфавит фиксирован).
// Для Unicode-строк нужно считать по рунам (map[rune]int) — здесь, как в
// оригинале задачи, латиница и байтового счётчика достаточно.
package main

func firstUniqChar(s string) int {
	var counts [256]int

	for i := 0; i < len(s); i++ {
		counts[s[i]]++
	}
	for i := 0; i < len(s); i++ {
		if counts[s[i]] == 1 {
			return i
		}
	}

	return -1
}

func main() {
	println("Test 1:", firstUniqChar("leetcode"))     // 0
	println("Test 2:", firstUniqChar("loveleetcode")) // 2
	println("Test 3:", firstUniqChar("aabb"))         // -1
}
