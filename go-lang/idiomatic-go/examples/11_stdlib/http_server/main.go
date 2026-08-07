// Минимальный сервер: собственный http.Server с обязательными таймаутами
// и обработчик — тип с методом ServeHTTP (интерфейс http.Handler).
package main

import (
	"net/http"
	"time"
)

type HelloHandler struct{}

func (hh HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// один метод — и тип уже http.Handler
	w.Write([]byte("Hello!\n"))
}

func main() {
	s := http.Server{
		Addr:        ":8080",
		ReadTimeout: 30 * time.Second,
		// таймауты обязательны: защита от зависших и вредоносных клиентов
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      HelloHandler{},
	}
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			// ErrServerClosed — штатное завершение, не ошибка
			panic(err)
		}
	}
}
