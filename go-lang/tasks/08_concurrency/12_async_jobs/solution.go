// MultiProcess: все JobFunc стартуют параллельно, первый успех возвращается,
// остальные отменяются; если упали все — возвращается последняя ошибка.
//
// Буферизованный канал результатов (len(jobs)) не даёт проигравшим
// горутинам зависнуть; производный context отменяет лишнюю работу.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Result struct{ Value string }

type JobFunc func(ctx context.Context, input string) (Result, error)

type jobResult struct {
	res Result
	err error
}

func MultiProcess(ctx context.Context, input string, jobs []JobFunc) (Result, error) {
	if len(jobs) == 0 {
		return Result{}, errors.New("no jobs")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // остановить оставшиеся задачи при выходе

	results := make(chan jobResult, len(jobs)) // буфер = нет утечек горутин

	for _, job := range jobs {
		go func(job JobFunc) {
			res, err := job(ctx, input)
			results <- jobResult{res: res, err: err}
		}(job)
	}

	var lastErr error
	for range jobs {
		select {
		case r := <-results:
			if r.err == nil {
				return r.res, nil // первый успех побеждает
			}
			lastErr = r.err
		case <-ctx.Done():
			return Result{}, ctx.Err()
		}
	}

	return Result{}, lastErr // упали все
}

func main() {
	slow := func(ctx context.Context, in string) (Result, error) {
		select {
		case <-time.After(300 * time.Millisecond):
			return Result{Value: "slow:" + in}, nil
		case <-ctx.Done():
			fmt.Println("slow: отменена")
			return Result{}, ctx.Err()
		}
	}
	fast := func(_ context.Context, in string) (Result, error) {
		time.Sleep(50 * time.Millisecond)
		return Result{Value: "fast:" + in}, nil
	}
	failing := func(_ context.Context, in string) (Result, error) {
		return Result{}, errors.New("boom")
	}

	res, err := MultiProcess(context.Background(), "task", []JobFunc{slow, fast, failing})
	fmt.Println(res, err) // {fast:task} <nil>

	res, err = MultiProcess(context.Background(), "task", []JobFunc{failing, failing})
	fmt.Println(res, err) // {} boom

	time.Sleep(100 * time.Millisecond) // дать slow напечатать отмену
}
