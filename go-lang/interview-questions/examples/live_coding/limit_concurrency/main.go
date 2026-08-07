// Задача: «как ограничишь количество параллельных запросов?»
// Ответ: семафор на буферизованном канале. Буфер = лимит.
//
// Два принципиально разных поведения при исчерпании лимита:
//  1. ЖДАТЬ свободный слот (обычная очередь) — semaphore;
//  2. сразу ОТКАЗАТЬ (противодавление, backpressure) — select с default.
//
// Третий вариант — ограничение по скорости, а не по параллелизму:
// time.Ticker / golang.org/x/time/rate (N запросов в секунду).
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var inFlight atomic.Int32 // счетчик, чтобы увидеть реальный параллелизм

func request(ctx context.Context, id int) (string, error) {
	cur := inFlight.Add(1)
	defer inFlight.Add(-1)
	fmt.Printf("  запрос %2d стартовал, сейчас в работе: %d\n", id, cur)

	select {
	case <-time.After(30 * time.Millisecond):
		return fmt.Sprintf("ответ %d", id), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// 1. Семафор: лишние горутины ЖДУТ свободный слот.
func withSemaphore(ctx context.Context, ids []int, limit int) []string {
	sem := make(chan struct{}, limit) // буфер = максимум одновременных запросов
	out := make([]string, len(ids))

	var wg sync.WaitGroup
	wg.Add(len(ids))
	for i, id := range ids {
		go func(i, id int) {
			defer wg.Done()

			select {
			case sem <- struct{}{}: // занять слот; если буфер полон — блокируемся здесь
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }() // освободить слот

			res, err := request(ctx, id)
			if err != nil {
				return
			}
			out[i] = res // каждый пишет в свой индекс — синхронизация не нужна
		}(i, id)
	}
	wg.Wait()
	return out
}

// 2. Противодавление: если слотов нет — не ждем, а сразу возвращаем ошибку.
type Limiter struct {
	sem chan struct{}
}

func NewLimiter(limit int) *Limiter { return &Limiter{sem: make(chan struct{}, limit)} }

var ErrBusy = errors.New("нет свободных слотов")

func (l *Limiter) Do(f func()) error {
	select {
	case l.sem <- struct{}{}:
		defer func() { <-l.sem }()
		f()
		return nil
	default:
		return ErrBusy // очередь не копится — клиент получит 429 и повторит позже
	}
}

func main() {
	ctx := context.Background()
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8}

	fmt.Println("1. Семафор, лимит 3 — лишние ждут:")
	res := withSemaphore(ctx, ids, 3)
	fmt.Println("  получено ответов:", len(res))

	fmt.Println("\n2. Противодавление, лимит 2 — лишним сразу отказ:")
	l := NewLimiter(2)
	var wg sync.WaitGroup
	var rejected atomic.Int32
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Do(func() { time.Sleep(20 * time.Millisecond) }); err != nil {
				rejected.Add(1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("  отказов:", rejected.Load())

	// Для HTTP-клиента есть и «бесплатный» способ ограничить параллелизм:
	// http.Transport{MaxConnsPerHost: N} — лимит на уровне соединений.
}
