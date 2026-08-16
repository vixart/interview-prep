// Контроль медленной зависимости: таймаут + fallback.
//
// Синхронный вызов getDiscount блокирует программу на неопределённое время.
// Оборачиваем вызов в горутину с буферизованным каналом и ждём результат
// через select с context.WithTimeout. Не успели — деградируем изящно:
// скидка 0, программа отвечает быстро.
package main

import (
	"context"
	"fmt"
	"time"
)

// Эта функция лезет по сети в старый монолит и может тупить.
func getDiscount() float64 {
	time.Sleep(2 * time.Second)
	return 12.0
}

func getDiscountCtx(ctx context.Context) (float64, error) {
	ch := make(chan float64, 1) // буфер: иначе горутина утечёт после таймаута

	go func() { ch <- getDiscount() }()

	select {
	case v := <-ch:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()
	discount, err := getDiscountCtx(ctx)
	if err != nil {
		discount = 0 // fallback: не тормозим программу из-за скидки
		fmt.Println("монолит не ответил вовремя:", err)
	}

	fmt.Printf("Ваша скидка: %v (ответ за %v)\n", discount, time.Since(start).Round(time.Millisecond))
	// Важно: сам getDiscount продолжает работать в фоне — реальную отмену
	// нужно пробрасывать внутрь (например, http.NewRequestWithContext).
}
