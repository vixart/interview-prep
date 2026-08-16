// Параллельные HTTP-запросы к списку URL (этапы 1-5 из условия).
//
// Итоговый вариант объединяет все этапы:
//   - запросы в горутинах, результаты — в канал (этап 2);
//   - количество URL заранее неизвестно: канал результатов закрывается
//     горутиной-закрывателем после wg.Wait (этап 3);
//   - отмена через context: после 2 успешных ответов остальные запросы
//     отменяются (этап 4);
//   - семафор ограничивает число одновременных запросов (этап 5).
//
// У клиента задан Timeout, тело ответа закрывается и вычитывается —
// иначе соединения не возвращаются в пул.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	url string
	ok  bool
	err error
}

func check(ctx context.Context, client *http.Client, url string) result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result{url: url, err: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return result{url: url, err: err}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body) // вычитать, чтобы соединение вернулось в пул

	return result{url: url, ok: resp.StatusCode == http.StatusOK}
}

func main() {
	urls := []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://eeeboy",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const maxParallel = 3
	sem := make(chan struct{}, maxParallel)

	results := make(chan result)
	var wg sync.WaitGroup
	var okCount atomic.Int64

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			sem <- struct{}{}        // занять слот
			defer func() { <-sem }() // освободить

			if ctx.Err() != nil {
				results <- result{url: url, err: ctx.Err()}
				return
			}

			r := check(ctx, client, url)
			if r.ok && okCount.Add(1) == 2 {
				cancel() // два успеха — отменяем остальные
			}
			results <- r
		}(url)
	}

	go func() { // закрыватель: сигнал "больше результатов не будет"
		wg.Wait()
		close(results)
	}()

	for r := range results { // печать в основном потоке
		switch {
		case r.ok:
			fmt.Println(r.url, "- ok")
		case r.err != nil:
			fmt.Println(r.url, "- not ok (", r.err, ")")
		default:
			fmt.Println(r.url, "- not ok")
		}
	}
}
