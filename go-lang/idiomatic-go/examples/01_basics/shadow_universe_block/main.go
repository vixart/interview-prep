// Затенить можно даже предопределенные идентификаторы из universe block.
// `true := 10` — легальный Go: true, nil, int, make не ключевые слова, а просто имена.
package main

import "fmt"

func main() {
	fmt.Println(true)
	true := 10
	// true — не ключевое слово, а имя из universe block, его можно затенить
	fmt.Println(true)
	// печатает 10, а не булево значение
}
