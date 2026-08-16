// Доработка SimulateRequest: потокобезопасность + таймаут + замер времени.
//
//  1. counter++ не атомарен (read-modify-write) -> atomic.AddInt64.
//  2. Таймаут: работа уходит в горутину, результат — в буферизованный канал
//     (буфер обязателен: иначе после таймаута горутина навсегда зависнет
//     на отправке — утечка). select ждёт результат или ctx.Done().
//  3. Время выполнения меряем time.Since.
package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var counter int64

func SimulateRequest(ctx context.Context) (int64, error) {
	done := make(chan int64, 1) // буфер: горутина не должна блокироваться

	go func() {
		time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
		done <- atomic.AddInt64(&counter, 1)
	}()

	select {
	case v := <-done:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			start := time.Now()
			val, err := SimulateRequest(ctx)
			if err != nil {
				log.Printf("запрос %d: ошибка %v (за %v)", i, err, time.Since(start))
				return
			}
			log.Printf("запрос %d: счетчик=%d (за %v)", i, val, time.Since(start))
		}(i)
	}

	wg.Wait()
}
