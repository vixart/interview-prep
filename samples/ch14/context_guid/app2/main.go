package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/learning-go-book-2e/ch14/sample_code/context_guid/tracker"
)

// Интерфейс бизнес-логики.
// ctx передаётся первым параметром, чтобы через него
// распространялись request-scoped данные (GUID), cancellation и дедлайны.
type Logic interface {
	QueryHandler(ctx context.Context, query string) (string, error)
}

type Controller struct {
	Logic Logic
}

func (c Controller) Second(rw http.ResponseWriter, req *http.Request) {

	// Каждый HTTP request уже содержит context.
	// Этот context был создан сервером и дополнен в middleware.
	ctx := req.Context()

	query := req.URL.Query().Get("query")

	// Передаём тот же context в бизнес-логику.
	// Благодаря этому GUID запроса доступен и в логике.
	result, err := c.Logic.QueryHandler(ctx, query)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write([]byte(result))
}

// Логгер принимает context,
// чтобы извлечь из него GUID запроса и добавить его в лог.
type Logger interface {
	Log(context.Context, string)
}

type LogicImpl struct {
	Logger Logger
	Remote string
}

func (l LogicImpl) QueryHandler(ctx context.Context, query string) (string, error) {

	// Логгер получает тот же context, который пришёл из HTTP запроса.
	// tracker.Logger извлекает GUID из context и добавляет его в лог.
	l.Logger.Log(ctx, "starting QueryHandler with query: "+query)

	return fmt.Sprintf("got query: '%s' from first", query), nil
}

func main() {

	r := chi.NewRouter()

	// Middleware выполняется для каждого входящего запроса.
	// Оно извлекает GUID из header (или генерирует новый)
	// и кладёт его в context:
	//
	// ctx = context.WithValue(ctx, guidKey, guid)
	// req = req.WithContext(ctx)
	r.Use(tracker.Middleware)

	controller := Controller{
		Logic: LogicImpl{
			Logger: tracker.Logger{},
		},
	}

	// endpoint второго сервиса
	// GET /second?query=...
	r.Get("/second", controller.Second)

	http.ListenAndServe(":4000", r)
}
