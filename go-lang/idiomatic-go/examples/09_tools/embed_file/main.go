// //go:embed встраивает файл в бинарник как строку (нужен импорт _ "embed").
// Файл passwords.txt после сборки не нужен: он уже внутри исполняемого файла.
// Запуск: go run . <пароль>
package main

import (
	_ "embed"
	// пустой импорт: сам пакет не используется, но без него //go:embed не работает
	"fmt"
	"os"
	"strings"
)

// Директива стоит ВПЛОТНУЮ к переменной — пустая строка между ними все сломает.
//
//go:embed passwords.txt
var passwords string // содержимое файла попадает сюда при сборке; в рантайме файл не нужен

// содержимое файла попадет сюда на этапе сборки; в рантайме файл уже не нужен

func main() {
	pwds := strings.Split(passwords, "\n")
	if len(os.Args) > 1 {
		for _, v := range pwds {
			if v == os.Args[1] {
				fmt.Println("true")
				os.Exit(0)
			}
		}
		fmt.Println("false")
	}
}
