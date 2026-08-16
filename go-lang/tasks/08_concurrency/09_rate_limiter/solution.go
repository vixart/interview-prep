// Rate limiter со скользящим окном ровно в 1 секунду.
//
// Требование «ни за какой промежуток длиной 1 секунда не более RPS» ломает
// fixed window (на стыке окон проходит до 2×RPS). Точный вариант —
// sliding window log: кольцевой буфер timestamp'ов последних RPS
// пропущенных запросов. Запрос пропускаем, если самый старый из них
// старше секунды. Память O(RPS) — при лимите 100k это ~800 КБ, приемлемо.
//
// Для экономии памяти на больших RPS используют sliding window counter
// или token bucket (golang.org/x/time/rate) — ценой небольшой неточности.
package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu    sync.Mutex
	times []time.Time // кольцевой буфер: моменты последних rps пропусков
	head  int
	size  int
	rps   int
	now   func() time.Time
}

func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{times: make([]time.Time, rps), rps: rps, now: time.Now}
}

// Allow возвращает true, если запрос можно пропустить.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()

	if rl.size == rl.rps {
		oldest := rl.times[rl.head]
		if now.Sub(oldest) < time.Second {
			return false // rps-й запрос назад был меньше секунды назад
		}
		rl.head = (rl.head + 1) % rl.rps // вытесняем самый старый
		rl.size--
	}

	rl.times[(rl.head+rl.size)%rl.rps] = now
	rl.size++
	return true
}

func main() {
	// Тест 1: не больше 3 запросов в любую секунду.
	rl := NewRateLimiter(3)
	allowed := 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			allowed++
		}
	}
	fmt.Println("мгновенно пропущено:", allowed) // 3

	// Тест 2: через секунду лимит снова доступен.
	time.Sleep(1100 * time.Millisecond)
	fmt.Println("после паузы:", rl.Allow()) // true

	// Тест 3: конкурентный доступ — суммарно не больше лимита.
	rl2 := NewRateLimiter(100)
	var passed int64
	var mu sync.Mutex
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				if rl2.Allow() {
					mu.Lock()
					passed++
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println("конкурентно пропущено:", passed, "(лимит 100)") // 100
}
