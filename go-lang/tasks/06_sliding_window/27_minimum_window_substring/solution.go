// Минимальное окно в s, содержащее все символы t (с учётом дубликатов).
//
// Окно переменного размера: расширяем правую границу, пока не набрали все
// символы (formed == required), затем сжимаем левую, фиксируя минимум.
// Время O(n + m), память O(алфавита).
package main

import "fmt"

func minWindow(s string, t string) string {
	if len(t) == 0 || len(s) < len(t) {
		return ""
	}

	var need, window [256]int
	required := 0 // сколько различных символов нужно набрать
	for i := 0; i < len(t); i++ {
		if need[t[i]] == 0 {
			required++
		}
		need[t[i]]++
	}

	formed := 0 // сколько символов уже набрано в нужном количестве
	bestLen, bestStart := -1, 0

	l := 0
	for r := 0; r < len(s); r++ {
		c := s[r]
		window[c]++
		if need[c] > 0 && window[c] == need[c] {
			formed++
		}

		for formed == required {
			if bestLen == -1 || r-l+1 < bestLen {
				bestLen, bestStart = r-l+1, l
			}

			lc := s[l]
			window[lc]--
			if need[lc] > 0 && window[lc] < need[lc] {
				formed--
			}
			l++
		}
	}

	if bestLen == -1 {
		return ""
	}
	return s[bestStart : bestStart+bestLen]
}

func main() {
	fmt.Println(minWindow("ADOBECODEBANC", "ABC")) // "BANC"
	fmt.Println(minWindow("a", "a"))               // "a"
	fmt.Println(minWindow("a", "aa"))              // ""
}
