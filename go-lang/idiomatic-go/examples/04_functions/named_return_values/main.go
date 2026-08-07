// Именованные возвращаемые значения: это заранее объявленные переменные.
// divAndRemainderConfusing показывает, почему они путают: присвоенные им значения
// могут быть проигнорированы явным return. Пустой return при них — источник багов.
package main

import (
	"errors"
	"fmt"
)

func divAndRemainder(numerator int, denominator int) (result int, remainder int,
	err error) {
	if denominator == 0 {
		err = errors.New("cannot divide by zero")
		// присваиваем именованному значению...
		return result, remainder, err
		// и возвращаем его явно — так читается лучше, чем голый return
	}
	result, remainder = numerator/denominator, numerator%denominator
	return result, remainder, err
}

func divAndRemainderConfusing(numerator, denominator int) (result int, remainder int,
	err error) {
	// assign some values
	result, remainder = 20, 30
	// ЛОВУШКА: эти значения ни на что не влияют — return ниже возвращает свои
	if denominator == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return numerator / denominator, numerator % denominator, nil
}

func main() {
	x, y, z := divAndRemainder(5, 2)
	fmt.Println(x, y, z)

	result, remainder, err := divAndRemainderConfusing(5, 2)
	fmt.Println(result, remainder, err)
}
