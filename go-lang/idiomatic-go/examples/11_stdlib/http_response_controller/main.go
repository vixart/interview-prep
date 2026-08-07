// Потоковый ответ: rc.Flush() отправляет уже записанную часть клиенту,
// не дожидаясь конца обработчика (HTTP-ответ идет chunked).
// Через ResponseController также доступны Hijack и дедлайны чтения/записи.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func doStuff(i int) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("%d\n", i*2)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	rc := http.NewResponseController(rw)
	// доступ к низкоуровневым возможностям ответа
	for i := 0; i < 10; i++ {
		result := doStuff(i)
		_, err := rw.Write([]byte(result))
		if err != nil {
			slog.Error("error writing", "msg", err)
			return
		}
		err = rc.Flush()
		// отправляем уже записанную часть клиенту, не дожидаясь конца обработчика
		if err != nil && !errors.Is(err, http.ErrNotSupported) {
			// не всякий ResponseWriter умеет Flush — это нормально
			slog.Error("error flushing", "msg", err)
			return
		}
	}
}

func main() {
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 6 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      http.HandlerFunc(handler),
	}
	err := s.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
