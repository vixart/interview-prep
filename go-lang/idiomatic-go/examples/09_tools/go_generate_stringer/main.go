// //go:generate: комментарий, по которому `go generate ./...` запускает stringer
// и создает метод String() для перечисления. Сгенерированный файл
// direction_string.go лежит рядом и коммитится в репозиторий.
package main

import "fmt"

type Direction int

const (
	_ Direction = iota
	North
	South
	East
	West
)

// Эту строку выполняет команда `go generate ./...` (сам компилятор ее игнорирует).
// Результат — файл direction_string.go рядом, его коммитят в репозиторий.
//
//go:generate stringer -type=Direction

func main() {
	fmt.Println(North.String())
	// метод String() пришел из сгенерированного файла, руками его не пишут
}
