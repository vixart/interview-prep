package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Задача: крутим цикл, пока:
// 1) не сработает timeout контекста
// 2) не выпадет число 1234
func main() {
	// создаём корневой context с дедлайном 2 секунды
	// после этого времени ctx.Done() будет закрыт
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)

	// важно: освобождаем ресурсы (таймер), даже если вышли раньше
	defer cancelFunc()

	total := 0
	count := 0

	for {
		select {
		case <-ctx.Done():
			// сюда попадаем, когда:
			// - истёк timeout (DeadlineExceeded)
			// - или кто-то вручную вызвал cancel
			// ctx.Err() объясняет причину
			fmt.Println("total:", total, "number of iterations:", count, ctx.Err())
			return

		default:
			// если контекст ещё "жив" — продолжаем работу
			// default нужен, чтобы select не блокировал цикл
		}

		newNum := rand.Intn(100_000_000)

		// альтернативное условие выхода (не связано с context)
		if newNum == 1_234 {
			fmt.Println("total:", total, "number of iterations:", count, "got 1,234")
			return
		}

		total += newNum
		count++
	}
}
