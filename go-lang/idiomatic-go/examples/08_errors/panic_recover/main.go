// recover работает только внутри defer: деление на 0 паникует,
// но функция восстанавливается и цикл продолжается.
// Так же прикрывают обработчик HTTP-запроса, чтобы одна паника не уронила процесс.
package main

import (
	"fmt"
)

func div60(i int) {
	defer func() {
		if v := recover(); v != nil {
			// recover работает ТОЛЬКО внутри defer; v — то, что передали в panic
			fmt.Println(v)
		}
	}()
	fmt.Println(60 / i)
	// на i == 0 здесь паника — но функция восстановится, а цикл продолжится
}

func main() {
	for _, val := range []int{1, 2, 0, 6} {
		div60(val)
	}
}
