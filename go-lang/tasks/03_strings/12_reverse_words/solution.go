// Реверс слов с сохранением палиндромов.
//
// Разбиваем строку на слова, каждое слово переворачиваем ПО РУНАМ
// (кириллица многобайтовая — по байтам переворачивать нельзя).
// Палиндром после переворота равен сам себе, поэтому отдельная проверка
// не обязательна, но оставлена явно — как в условии.
package main

import (
	"fmt"
	"strings"
)

func reverseWords(s string) string {
	words := strings.Split(s, " ")

	for i, w := range words {
		if !isPalindrome(w) {
			words[i] = reverseRunes(w)
		}
	}

	return strings.Join(words, " ")
}

func reverseRunes(s string) string {
	r := []rune(s)
	for l, rr := 0, len(r)-1; l < rr; l, rr = l+1, rr-1 {
		r[l], r[rr] = r[rr], r[l]
	}
	return string(r)
}

func isPalindrome(s string) bool {
	r := []rune(s)
	for l, rr := 0, len(r)-1; l < rr; l, rr = l+1, rr-1 {
		if r[l] != r[rr] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(reverseWords("Hello worlD ollo")) // "olleH Dlrow ollo"
	fmt.Println(reverseWords("привет мир ара"))   // "тевирп рим ара"
}
