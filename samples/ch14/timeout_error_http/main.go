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
// Если context отменится (cancel/timeout), http.Client прервёт запрос.
func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func main() {

	// Создаём родительский context с таймаутом 3 секунды.
	// После истечения таймаута ctx.Done() закроется автоматически.
	ctx, cancelFuncParent := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFuncParent()

	// Поверх него создаём context с возможностью указать причину отмены.
	// Теперь мы можем не только отменить context,
	// но и передать конкретную ошибку (cause).
	ctx, cancelFunc := context.WithCancelCause(ctx)
	defer cancelFunc(nil)

	ch := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine №1 — делает периодические HTTP запросы.
	go func() {
		defer wg.Done()

		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/status/200,200,200,500")

			// Если context отменён (например таймаутом),
			// HTTP клиент вернёт ошибку.
			if err != nil {
				cancelFunc(fmt.Errorf("in status goroutine: %w", err))
				return
			}

			// Если сервер вернул 500 — отменяем общий context.
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

			// Если context уже отменён (cancel или timeout),
			// запрос завершится ошибкой.
			if err != nil {
				fmt.Println("in delay goroutine:", err)

				cancelFunc(fmt.Errorf("in delay goroutine: %w", err))
				return
			}

			ch <- "success from delay: " + resp.Header.Get("date")
		}
	}()

loop:
	for {
		select {

		// Получаем результаты от goroutine
		case s := <-ch:
			fmt.Println("in main:", s)

		// ctx.Done() закрывается когда:
		// 1) истёк таймаут WithTimeout
		// 2) была вызвана cancelFunc(...)
		case <-ctx.Done():

			// context.Cause(ctx) — причина отмены (если была передана)
			// ctx.Err() — тип отмены (context.Canceled или context.DeadlineExceeded)
			fmt.Println(
				"in main: cancelled with cause:",
				context.Cause(ctx),
				"err:",
				ctx.Err(),
			)

			break loop
		}
	}

	// Ждём завершения goroutine.
	// Они завершатся, потому что используют тот же context.
	wg.Wait()

	fmt.Println("context cause:", context.Cause(ctx))
}
