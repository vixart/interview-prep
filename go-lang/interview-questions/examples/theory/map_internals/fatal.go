//go:build mapfatal

// Демонстрация того, что конкурентная запись в map — это FATAL ERROR, а не паника:
// ее нельзя перехватить через recover, процесс просто умирает.
// Собирается только с тегом, чтобы обычный make check оставался зеленым:
//
//	go run -tags mapfatal ./theory/map_internals
package main

import (
	"fmt"
	"sync"
)

func concurrentDemo() {
	defer func() {
		// recover НЕ сработает: fatal error не является паникой
		if r := recover(); r != nil {
			fmt.Println("перехвачено:", r)
		}
	}()

	m := map[int]int{}
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100_000; j++ {
				m[base+j] = j // fatal error: concurrent map writes
			}
		}(i * 1_000_000)
	}
	wg.Wait()
	fmt.Println("сюда не дойдем:", len(m))
}
