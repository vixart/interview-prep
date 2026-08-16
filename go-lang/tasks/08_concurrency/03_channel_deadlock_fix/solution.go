// Исправленный вариант задачи про тикер и канал.
//
// Проблемы оригинала:
//   - main писала в небуферизованный канал внутри select и сама же была
//     единственным читателем -> deadlock;
//   - ticker не останавливался (утечка);
//   - через 3 секунды горутина писала в канал, который уже никто не читал.
//
// Решение: канал только для сигнала от горутины, тикер — отдельная ветка
// select, defer ticker.Stop(), выход по получению сигнала.
package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool, 1)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Отдельная горутина отвисла")
		done <- false
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Произошёл тик тикера")
		case value := <-done:
			fmt.Printf("Получено значение %t\n", value)
			return
		}
	}
}
