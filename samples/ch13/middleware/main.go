package main

import (
	"log/slog"
	"net/http"
	"time"
)

func main() {
	// TerribleSecurityProvider — функция-генератор middleware.
	// Она принимает параметр (пароль) и возвращает middleware:
	// func(http.Handler) http.Handler
	terribleSecurity := TerribleSecurityProvider("GOPHER")

	mux := http.NewServeMux()

	// Здесь создаётся цепочка middleware вокруг handler.
	//
	// Порядок оборачивания:
	//   terribleSecurity(RequestTimer(handler))
	//
	// Но выполняться они будут так:
	//   request
	//     ↓
	//   terribleSecurity
	//     ↓
	//   RequestTimer
	//     ↓
	//   handler
	//
	// middleware вызывают следующий handler через h.ServeHTTP.
	mux.Handle("/hello", terribleSecurity(RequestTimer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello!\n"))
		}),
	)))

	// Альтернативный способ использования middleware:
	// можно обернуть весь router, потому что ServeMux
	// тоже реализует интерфейс http.Handler.
	//
	// mux = terribleSecurity(RequestTimer(mux))
	//
	// Тогда middleware будут применяться ко всем маршрутам.

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

// RequestTimer — пример middleware.
// Принимает handler и возвращает новый handler,
// который выполняет код до и после вызова следующего обработчика.
func RequestTimer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// код "до" вызова следующего обработчика
		start := time.Now()

		// передача управления следующему handler в цепочке
		h.ServeHTTP(w, r)

		// код "после" выполнения handler
		slog.Info(
			"request time",
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

var securityMsg = []byte("You didn't give the secret password\n")

// TerribleSecurityProvider — middleware factory.
// Сначала вызывается эта функция (с параметрами конфигурации),
// а она возвращает middleware.
func TerribleSecurityProvider(password string) func(http.Handler) http.Handler {

	// Возвращаем middleware.
	return func(h http.Handler) http.Handler {

		// Middleware — это новый handler,
		// который решает, вызывать ли следующий handler.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Если проверка не прошла — цепочка middleware
			// прерывается (следующий handler не вызывается).
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(securityMsg)
				return
			}

			// Если проверка прошла —
			// управление передается следующему handler.
			h.ServeHTTP(w, r)
		})
	}
}
