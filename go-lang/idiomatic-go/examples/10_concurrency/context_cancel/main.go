// Остановка горутины контекстом: генератор пишет в канал через select с ctx.Done(),
// поэтому после break в main и вызова cancel() он не остается заблокированным навсегда.
// Сравни с goroutine_leak, где такой защиты нет.
package main

import (
	"context"
	"fmt"
)

func countTo(ctx context.Context, max int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		// горутина сама закрывает канал при выходе
		for i := 0; i < max; i++ {
			select {
			case <-ctx.Done():
				// выход по отмене — иначе горутина зависла бы на записи в ch
				return
			case ch <- i:
				// select ждет: либо отдали значение, либо отменили
			}
		}
	}()
	return ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := countTo(ctx, 10)
	for i := range ch {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
	cancel()
	// без этого вызова горутина осталась бы заблокированной навсегда
}
