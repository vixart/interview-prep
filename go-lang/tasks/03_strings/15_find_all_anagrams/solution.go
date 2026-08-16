// Поиск всех анаграмм строки p в строке s — скользящее окно фиксированной
// длины len(p) с частотными счётчиками.
//
// Держим счётчик частот окна и сравниваем его с частотами p не «в лоб»,
// а через счётчик matched — сколько символов уже совпало по количеству.
// Время O(n), память O(1) (алфавит фиксирован — 256 байтовых значений).
package main

import "fmt"

func findAnagrams(s string, p string) []int {
	if len(p) > len(s) {
		return nil
	}

	var need, window [256]int
	for i := 0; i < len(p); i++ {
		need[p[i]]++
	}

	result := []int{}
	matched := 0
	distinct := 0
	for _, c := range need {
		if c > 0 {
			distinct++
		}
	}

	add := func(c byte) {
		window[c]++
		if window[c] == need[c] {
			matched++
		} else if window[c] == need[c]+1 {
			matched--
		}
	}
	del := func(c byte) {
		if window[c] == need[c] {
			matched--
		} else if window[c] == need[c]+1 {
			matched++
		}
		window[c]--
	}

	for i := 0; i < len(s); i++ {
		add(s[i])
		if i >= len(p) {
			del(s[i-len(p)])
		}
		if i >= len(p)-1 && matched == distinct {
			result = append(result, i-len(p)+1)
		}
	}

	return result
}

func main() {
	fmt.Println(findAnagrams("cbaebabacd", "abc")) // [0 6]
	fmt.Println(findAnagrams("abab", "ab"))        // [0 1 2]
}
