//go:build copylocks

// Демонстрация того, за что ругается go vet. Файл собирается только с тегом,
// чтобы обычный `make check` оставался зеленым:
//
//	go vet -tags copylocks ./problems/mutex_copy   # увидишь "passes lock by value"
//	go run -tags copylocks -race ./problems/mutex_copy
package main

import (
	"fmt"
	"sync"
)

// ПРИЕМНИК ПО ЗНАЧЕНИЮ — каждый вызов работает с КОПИЕЙ мьютекса и поля n.
func (c Counter) IncBroken() {
	c.mu.Lock() // блокируем копию: другие горутины ее не видят
	defer c.mu.Unlock()
	c.n++ // инкрементим копию: наружу изменение не попадет
}

// Передача структуры с мьютексом по значению — тот же дефект.
func incAll(c Counter, times int) {
	for i := 0; i < times; i++ {
		c.IncBroken()
	}
}

func runBadCopy() {
	var c Counter
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			c.IncBroken() // гонка + потерянные инкременты одновременно
		}()
	}
	wg.Wait()
	fmt.Println("сломано, приемник по значению:", c.n) // 0

	incAll(c, 10)
	fmt.Println("после incAll по значению:", c.n) // тоже 0

	// Отдельный сценарий: копирование ЗАЛОЧЕННОГО мьютекса.
	var mu sync.Mutex
	mu.Lock()
	cp := mu    // копия сохраняет состояние "залочен"
	mu.Unlock() // разблокировали оригинал, но не копию
	if cp.TryLock() {
		fmt.Println("копия свободна — так не будет")
	} else {
		fmt.Println("копия осталась залоченной: обычный Lock здесь завис бы навсегда")
	}
}
