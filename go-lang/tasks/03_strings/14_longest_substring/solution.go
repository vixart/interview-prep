// Самая длинная подстрока без повторяющихся символов.
//
// Скользящее окно [left, i] + мапа "символ → последний индекс".
// Встретили символ, который уже есть в окне, — передвигаем left
// за его предыдущее вхождение. Время O(n), память O(k).
package main

import "fmt"

func lengthOfLongestSubstring(s string) int {
	lastSeen := make(map[byte]int)
	best, left := 0, 0

	for i := 0; i < len(s); i++ {
		if pos, ok := lastSeen[s[i]]; ok && pos >= left {
			left = pos + 1
		}
		lastSeen[s[i]] = i

		if i-left+1 > best {
			best = i - left + 1
		}
	}

	return best
}

func main() {
	fmt.Println(lengthOfLongestSubstring("abcabcbb")) // 3 ("abc")
	fmt.Println(lengthOfLongestSubstring("bbbbb"))    // 1 ("b")
	fmt.Println(lengthOfLongestSubstring("pwwkew"))   // 3 ("wke")
}
