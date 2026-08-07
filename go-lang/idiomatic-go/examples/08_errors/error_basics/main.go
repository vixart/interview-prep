// Базовое соглашение: ошибка — последнее возвращаемое значение, при ошибке
// остальные значения нулевые, вызывающий проверяет `if err != nil`.
package main

import (
	"errors"
	"fmt"
	"os"
)

func calcRemainderAndMod(numerator, denominator int) (int, int, error) {
	if denominator == 0 {
		return 0, 0, errors.New("denominator is 0")
		// при ошибке остальным значениям — нулевые значения
	}
	return numerator / denominator, numerator % denominator, nil
}

func main() {
	numerator := 20
	denominator := 3
	remainder, mod, err := calcRemainderAndMod(numerator, denominator)
	if err != nil {
		// проверка ошибки сразу после вызова — основной ритм Go-кода
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(remainder, mod)
}
