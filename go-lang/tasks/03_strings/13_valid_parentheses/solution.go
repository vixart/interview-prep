// Валидация скобочной последовательности через стек.
//
// Открывающие скобки кладём в стек; для закрывающей проверяем,
// что на вершине стека лежит парная открывающая. В конце стек должен
// быть пуст. Время O(n), память O(n).
package main

import "fmt"

func isValid(s string) bool {
	pairs := map[byte]byte{')': '(', ']': '[', '}': '{'}
	var stack []byte

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {
		case '(', '[', '{':
			stack = append(stack, c)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[c] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func main() {
	fmt.Println(isValid("()"))     // true
	fmt.Println(isValid("()[]{}")) // true
	fmt.Println(isValid("([])"))   // true
	fmt.Println(isValid("([]}"))   // false (неправильный тип)
	fmt.Println(isValid("({)}"))   // false (неправильный порядок)
}
