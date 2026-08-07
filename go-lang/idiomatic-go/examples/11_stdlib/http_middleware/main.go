// Middleware — функция func(http.Handler) http.Handler.
// RequestTimer замеряет время, TerribleSecurityProvider проверяет заголовок и может
// не вызвать следующий обработчик. Оборачивать можно как один маршрут, так и весь mux.
package main

import (
	"log/slog"
	"net/http"
	"time"
)

func main() {
	terribleSecurity := TerribleSecurityProvider("GOPHER")

	mux := http.NewServeMux()

	// to apply the middleware to just the single route
	mux.Handle("/hello", terribleSecurity(RequestTimer(
		// обертки применяются снаружи внутрь: security → timer → handler
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello!\n"))
		}))))

	// or to apply the middleware to every route in the mux:
	//
	//	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("Hello!\n"))
	//	})
	//	mux = terribleSecurity(RequestTimer(mux))

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

func RequestTimer(h http.Handler) http.Handler {
	// сигнатура middleware: принимает Handler, возвращает Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		// вызов следующего обработчика — до и после можно вставить свою логику
		end := time.Now()
		slog.Info("request time", "path", r.URL.Path, "duration", end.Sub(start))
	})
}

var securityMsg = []byte("You didn't give the secret password\n")

func TerribleSecurityProvider(password string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(securityMsg)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
