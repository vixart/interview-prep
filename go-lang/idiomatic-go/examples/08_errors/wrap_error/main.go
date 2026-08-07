// Обертывание ошибки: %w сохраняет исходную ошибку в цепочке,
// errors.Unwrap достает ее обратно. Вывод — две строки: с контекстом и без.
package main

import (
	"errors"
	"fmt"
	"os"
)

func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("in fileChecker: %w", err)
		// %w сохраняет исходную ошибку в цепочке (%v сохранил бы только текст)
	}
	f.Close()
	return nil
}

func main() {
	err := fileChecker("not_here.txt")
	if err != nil {
		fmt.Println(err)
		if wrappedErr := errors.Unwrap(err); wrappedErr != nil {
			// Unwrap достает обернутую ошибку; на практике вместо него берут Is/As
			fmt.Println(wrappedErr)
		}
	}
}
