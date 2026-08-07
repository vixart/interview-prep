// Строка → []byte (кодовые единицы UTF-8) и строка → []rune (кодовые точки).
// Смотри на длину вывода: эмодзи занимает 4 байта, но всего одну руну.
package main

import "fmt"

func main() {
	var s string = "Hello, 🌞"
	var bs []byte = []byte(s)
	// байты UTF-8: эмодзи займет 4 элемента
	var rs []rune = []rune(s)
	// кодовые точки: тот же эмодзи — ровно 1 элемент
	fmt.Println(bs)
	fmt.Println(rs)
}
