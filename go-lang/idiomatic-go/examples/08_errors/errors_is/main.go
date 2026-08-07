// errors.Is ищет конкретную ошибку по ВСЕЙ цепочке:
// os.ErrNotExist находится, хотя ошибка была обернута через %w.
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
	}
	f.Close()
	return nil
}

func main() {
	err := fileChecker("not_here.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// ошибка обернута через %w, но Is находит ее по всей цепочке
			fmt.Println("That file doesn't exist")
		}
	}
}
