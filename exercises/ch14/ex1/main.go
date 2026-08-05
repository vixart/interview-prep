package main

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

// Timeout добавляет deadline в context каждого запроса.
// Важно: context не "прерывает" выполнение сам по себе —
// он лишь посылает сигнал отмены (через Done),
// а код ниже должен уметь его корректно обработать.
func Timeout(ms int) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// создаём новый context с дедлайном
			// этот ctx "оборачивает" родительский и унаследует его отмену
			ctx, cancelFunc := context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)

			// обязательно освобождаем ресурсы (таймер, ссылки)
			defer cancelFunc()

			// кладём новый context в request
			// дальше ВСЯ цепочка должна использовать именно его
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

func main() {
	middleware := Timeout(100)

	server := http.Server{
		// теперь каждый запрос будет с timeout-контекстом
		Handler: middleware(http.HandlerFunc(sleepy)),
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

func sleepy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// передаём context дальше — ключевая идея Go:
	// все долгие/внешние операции должны его принимать
	message, err := doThing(ctx)

	if err != nil {
		// если deadline exceeded — значит ctx отменён по таймауту
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write([]byte(message))
}

func doThing(ctx context.Context) (string, error) {
	wait := rand.Intn(200)

	select {
	case <-time.After(time.Duration(wait) * time.Millisecond):
		// "работа" завершилась раньше дедлайна
		return "Done!", nil

	case <-ctx.Done():
		// ключевой момент:
		// ctx.Done() закрывается, когда:
		// - вышел timeout
		// - или родительский ctx отменён
		// ctx.Err() скажет причину (DeadlineExceeded / Canceled)
		return "Too slow!", ctx.Err()
	}
}
