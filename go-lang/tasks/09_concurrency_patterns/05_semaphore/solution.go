// Semaphore: ограничение числа одновременных загрузок.
//
// Семафор на буферизованном канале: Acquire кладёт токен (блокируется,
// когда буфер полон), Release забирает. В отличие от worker pool здесь
// горутина на задачу, а ограничивается только одновременность.
// AcquireCtx позволяет отказаться от ожидания по таймауту/отмене.
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore { return make(Semaphore, n) }

func (s Semaphore) Acquire() { s <- struct{}{} }
func (s Semaphore) Release() { <-s }

func (s Semaphore) AcquireCtx(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	const limit = 3
	sem := NewSemaphore(limit)

	var active, peak atomic.Int64
	download := func(url string) {
		cur := active.Add(1)
		defer active.Add(-1)
		for { // зафиксировать пик одновременности
			p := peak.Load()
			if cur <= p || peak.CompareAndSwap(p, cur) {
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
		fmt.Println("загружен", url)
	}

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release() // освобождение даже при панике

			download(fmt.Sprintf("file-%d", i))
		}(i)
	}
	wg.Wait()

	fmt.Println("пик одновременных загрузок:", peak.Load(), "(лимит", limit, ")")
}
