// Задача: «есть несколько источников данных — как объединишь?»
// Ответ: fan-in — по горутине на источник, все пишут в один общий канал,
// закрывает его отдельная горутина после WaitGroup.
//
// Три момента, которые ждут от тебя на собеседовании:
//  1. общий канал закрывается ОДИН раз и только после всех писателей;
//  2. отмена по контексту — иначе горутины повиснут, если читатель ушел;
//  3. порядок значений не гарантирован (нужен порядок — нумеруй элементы).
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// source имитирует источник: канал, который сам закроется, отдав n значений.
func source(ctx context.Context, name string, n int, delay time.Duration) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out) // закрывает тот, кто пишет
		for i := 1; i <= n; i++ {
			select {
			case out <- fmt.Sprintf("%s-%d", name, i):
				time.Sleep(delay)
			case <-ctx.Done():
				return // читатель ушел — не зависаем на записи
			}
		}
	}()
	return out
}

// FanIn сливает несколько каналов в один.
func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T {
	out := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch { // цикл закончится, когда источник закроется
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
		// единственное место, где закрывается общий канал:
		// если бы это делал каждый писатель — была бы паника «close of closed channel»
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := source(ctx, "db", 3, 15*time.Millisecond)
	b := source(ctx, "cache", 4, 10*time.Millisecond)
	c := source(ctx, "api", 2, 25*time.Millisecond)

	for v := range FanIn(ctx, a, b, c) {
		// значения приходят вперемешку — это нормально для fan-in
		fmt.Println(v)
	}
	fmt.Println("все источники исчерпаны")

	// Варианты вопроса:
	// - «а если один источник вечно висит?» → context.WithTimeout;
	// - «а если нужен только первый ответ?» → select по каналам + cancel остальным;
	// - «а если источников тысячи?» → ограничить параллелизм (см. limit_concurrency).
}
