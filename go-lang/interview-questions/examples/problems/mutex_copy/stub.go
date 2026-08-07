//go:build !copylocks

package main

import "fmt"

// Заглушка для сборки без тега copylocks: сам «плохой» код лежит в badcopy.go.
func runBadCopy() {
	fmt.Println("демонстрация копирования мьютекса: go run -tags copylocks ./problems/mutex_copy")
}
