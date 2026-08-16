// Параллельный поиск на нескольких серверах: первый успех побеждает.
//
//   - горутина на сервер, результаты в БУФЕРИЗОВАННЫЙ канал (len(servers)),
//     чтобы проигравшие горутины не зависли на отправке (утечка);
//   - первый успешный ответ возвращается сразу, остальные отменяются
//     через context;
//   - если все вернули ошибку — отдаём агрегированную ошибку (errors.Join).
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type SearchFunc func(server, query string) ([]string, error)

type searchResult struct {
	res []string
	err error
}

func Search(servers []string, query string, searchFunc SearchFunc) ([]string, error) {
	if len(servers) == 0 {
		return nil, errors.New("no servers")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // отменить оставшиеся запросы при выходе

	results := make(chan searchResult, len(servers)) // буфер = нет утечек

	for _, server := range servers {
		go func(server string) {
			res, err := searchFunc(server, query)
			if ctx.Err() != nil {
				return // победитель уже найден
			}
			results <- searchResult{res: res, err: err}
		}(server)
	}

	var errs []error
	for range servers {
		r := <-results
		if r.err == nil {
			return r.res, nil // первый успех
		}
		errs = append(errs, r.err)
	}

	return nil, fmt.Errorf("all servers failed: %w", errors.Join(errs...))
}

func main() {
	demoSearch := func(server, query string) ([]string, error) {
		switch {
		case strings.Contains(server, "slow"):
			time.Sleep(300 * time.Millisecond)
			return []string{query + "@" + server}, nil
		case strings.Contains(server, "bad"):
			return nil, fmt.Errorf("%s: 500", server)
		default:
			time.Sleep(50 * time.Millisecond)
			return []string{query + "@" + server}, nil
		}
	}

	res, err := Search([]string{"bad-1", "slow-1", "fast-1"}, "golang", demoSearch)
	fmt.Println(res, err) // [golang@fast-1] <nil>

	res, err = Search([]string{"bad-1", "bad-2"}, "golang", demoSearch)
	fmt.Println(res, err) // [] all servers failed: ...
}
