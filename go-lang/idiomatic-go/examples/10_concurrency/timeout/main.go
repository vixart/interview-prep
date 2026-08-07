// Таймаут на произвольную работу: worker в отдельной горутине, select между результатом
// и ctx.Done(). Ключевая деталь — буфер 1 у канала out: если таймаут сработал первым,
// горутина все равно запишет результат и завершится, а не утечет.
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	result, err := timeLimit(doSomeWork, 2*time.Second)
	fmt.Println(result, err)
}

func timeLimit[T any](worker func() T, limit time.Duration) (T, error) {
	out := make(chan T, 1)
	// БУФЕР 1 — критично: при таймауте горутина запишет результат и завершится, а не утечет
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	// освобождает таймер контекста, даже если worker успел раньше
	go func() {
		out <- worker()
	}()
	select {
	case result := <-out:
		return result, nil
	case <-ctx.Done():
		// канал закрывается по истечении limit
		var zero T
		return zero, errors.New("work timed out")
	}
}

func doSomeWork() int {
	if x := rand.Int(); x%2 == 0 {
		return x
	} else {
		time.Sleep(10 * time.Second)
		return 100
	}
}
