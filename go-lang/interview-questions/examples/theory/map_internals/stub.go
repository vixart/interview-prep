//go:build !mapfatal

package main

import "fmt"

// Заглушка: сама демонстрация конкурентной записи лежит в fatal.go за тегом сборки.
func concurrentDemo() {
	fmt.Println("5. конкурентная запись в map: go run -tags mapfatal ./theory/map_internals")
}
