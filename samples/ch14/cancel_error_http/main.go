package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// makeRequest создаёт HTTP запрос, связанный с context.
// Если ctx будет отменён, http.Client автоматически
// прервёт выполнение запроса.
func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func main() {

	// Создаём корневой context для всей программы.
	// WithCancelCause позволяет отменить context
	// и указать причину отмены (ошибку).
	ctx, cancelFunc := context.WithCancelCause(context.Background())

	// На случай нормального завершения программы
	// вызываем cancel без причины.
	defer cancelFunc(nil)

	ch := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine №1 — периодически делает HTTP запрос.
	go func() {
		defer wg.Done()

		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/status/200,200,200,500")

			// Если запрос завершился ошибкой,
			// отменяем общий context для всех goroutine.
			if err != nil {
				cancelFunc(fmt.Errorf("in status goroutine: %w", err))
				return
			}

			// Если сервер вернул 500 — тоже отменяем context.
			if resp.StatusCode == http.StatusInternalServerError {
				cancelFunc(errors.New("bad status"))
				return
			}

			ch <- "success from status"

			time.Sleep(1 * time.Second)
		}
	}()

	// Goroutine №2 — делает медленный HTTP запрос.
	go func() {
		defer wg.Done()

		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/delay/1")

			// Если context уже отменён,
			// HTTP клиент вернёт ошибку.
			if err != nil {
				fmt.Println("in delay goroutine:", err)

				// Пробуем отменить context с причиной
				cancelFunc(fmt.Errorf("in delay goroutine: %w", err))
				return
			}

			ch <- "success from delay: " + resp.Header.Get("date")
		}
	}()

loop:
	for {
		select {

		// Получаем сообщения от goroutine
		case s := <-ch:
			fmt.Println("in main:", s)

		// ctx.Done() закрывается, когда context отменяется.
		// Это основной механизм сигнализации отмены.
		case <-ctx.Done():

			// context.Cause возвращает причину отмены,
			// переданную в cancelFunc.
			fmt.Println("in main: cancelled with error", context.Cause(ctx))

			break loop
		}
	}

	// Ждём завершения всех goroutine.
	// Они завершатся потому что используют тот же context.
	wg.Wait()

	// Можно ещё раз получить причину отмены context.
	fmt.Println("context cause:", context.Cause(ctx))
}
