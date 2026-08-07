// Утечка горутины и ее лечение контекстом.
// Правило: запуская горутину, сразу решай, как она завершится.
package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// leaky: если читатель прекратит читать раньше времени, горутина навсегда
// зависнет на "ch <- i" — планировщик продолжит держать ее в памяти.
func leaky(max int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < max; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

// fixed: select с ctx.Done() дает горутине выход в любой момент.
func fixed(ctx context.Context, max int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < max; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

func main() {
	start := runtime.NumGoroutine()

	for i := range leaky(1_000) {
		if i > 2 {
			break // горутина остается заблокированной навсегда
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	for i := range fixed(ctx, 1_000) {
		if i > 2 {
			break
		}
	}
	cancel() // без этого вызова вторая горутина утекла бы точно так же

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("горутин было %d, стало %d (одна утекла)\n", start, runtime.NumGoroutine())
}
