package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

/*
Сервер возвращает текущее время.

Формат ответа выбирается через HTTP заголовок Accept:

Accept: application/json → JSON
любой другой / отсутствует → текст (RFC3339)
*/

func main() {

	// Настройка handler для slog.
	// Handler определяет формат логов и место вывода.
	options := &slog.HandlerOptions{}

	// JSON формат логов, вывод в stderr.
	handler := slog.NewJSONHandler(os.Stderr, options)

	// Создание собственного logger.
	mySlog := slog.New(handler)

	// Создание HTTP router с middleware логирования.
	r := createChiRouter(mySlog)

	// Конфигурация HTTP сервера.
	s := http.Server{
		Addr:         ":8080",           // порт сервера
		ReadTimeout:  30 * time.Second,  // таймаут чтения запроса
		WriteTimeout: 90 * time.Second,  // таймаут записи ответа
		IdleTimeout:  120 * time.Second, // таймаут keep-alive соединений
		Handler:      r,                 // router как обработчик
	}

	// Запуск HTTP сервера.
	err := s.ListenAndServe()

	// Игнорируем штатное завершение сервера.
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func createChiRouter(logger *slog.Logger) chi.Router {

	// Создание router и добавление middleware через With().
	r := chi.NewRouter().With(func(handler http.Handler) http.Handler {

		// Middleware оборачивает каждый запрос.
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

			// Из RemoteAddr извлекается IP клиента.
			ip, _, _ := strings.Cut(req.RemoteAddr, ":")

			// Логируем IP входящего запроса.
			logger.Info("incoming IP", "ip", ip)

			// Передаем управление следующему handler.
			handler.ServeHTTP(rw, req)
		})
	})

	// Обработчик GET /
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		// Текущее время.
		now := time.Now()

		var out string

		// Content negotiation через Accept header.
		if r.Header.Get("Accept") == "application/json" {

			// Формирование JSON ответа.
			out = buildJSON(now)

		} else {

			// Текстовый ответ.
			out = buildText(now)
		}

		// HTTP статус ответа.
		w.WriteHeader(http.StatusOK)

		// Отправка тела ответа.
		w.Write([]byte(out))
	})

	return r
}

// Возвращает время в текстовом формате RFC3339.
func buildText(now time.Time) string {
	return now.Format(time.RFC3339)
}

/*
JSON формат ответа:

{
  "day_of_week": "Monday",
  "day_of_month": 10,
  "month": "April",
  "year": 2023,
  "hour": 20,
  "minute": 15,
  "second": 20
}
*/

func buildJSON(now time.Time) string {

	// Анонимная структура используется
	// как DTO для сериализации JSON.
	timeOut := struct {
		DayOfWeek  string `json:"day_of_week"`
		DayOfMonth int    `json:"day_of_month"`
		Month      string `json:"month"`
		Year       int    `json:"year"`
		Hour       int    `json:"hour"`
		Minute     int    `json:"minute"`
		Second     int    `json:"second"`
	}{
		// Заполнение структуры из объекта time.Time.
		DayOfWeek:  now.Weekday().String(),
		DayOfMonth: now.Day(),
		Month:      now.Month().String(),
		Year:       now.Year(),
		Hour:       now.Hour(),
		Minute:     now.Minute(),
		Second:     now.Second(),
	}

	// Сериализация структуры в JSON.
	out, _ := json.Marshal(timeOut)

	// Возврат JSON как строки.
	return string(out)
}
