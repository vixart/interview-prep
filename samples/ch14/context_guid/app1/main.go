package main

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/learning-go-book-2e/ch14/sample_code/context_guid/tracker"
)

// Logic — интерфейс бизнес-логики.
// ctx передаётся первым аргументом, чтобы через него
// распространялись cancellation, deadlines и request-scoped данные (GUID).
type Logic interface {
	Process(ctx context.Context, data string) (string, error)
}

type Controller struct {
	Logic Logic
}

func (c Controller) First(rw http.ResponseWriter, req *http.Request) {

	// Каждый HTTP request уже содержит context.
	// Этот context был создан сервером и дополнен в middleware.
	ctx := req.Context()

	data := req.URL.Query().Get("data")

	// Передаём context дальше в бизнес-логику.
	// Таким образом весь стек вызовов работает в рамках одного request context.
	result, err := c.Logic.Process(ctx, data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write([]byte(result))
}

// Logger принимает context,
// чтобы извлечь из него GUID запроса и добавить его в лог.
type Logger interface {
	Log(context.Context, string)
}

// RequestDecorator — функция, которая модифицирует исходящий HTTP request.
// В нашем случае она копирует GUID из context в header.
type RequestDecorator func(*http.Request) *http.Request

type LogicImpl struct {
	RequestDecorator RequestDecorator
	Logger           Logger
	Remote           string
}

func (l LogicImpl) Process(ctx context.Context, data string) (string, error) {

	// Логгер получает тот же context,
	// поэтому может извлечь GUID запроса.
	l.Logger.Log(ctx, "starting Process with "+data)

	// Создаём исходящий HTTP запрос.
	// Важно: используем NewRequestWithContext,
	// чтобы тот же context передался в HTTP client.
	//
	// Это позволяет:
	//   - отменить HTTP вызов если клиент отменил запрос
	//   - передать дедлайн
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		l.Remote+"/second?query="+data,
		nil,
	)

	if err != nil {
		l.Logger.Log(ctx, "error building remote request:"+err.Error())
		return "", err
	}

	// Декоратор добавляет GUID из context в header запроса.
	// Это позволяет передать correlation id в другой сервис.
	req = l.RequestDecorator(req)

	// Выполняем HTTP запрос.
	// Если context отменится — client тоже отменит запрос.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Logger.Log(ctx, "error building remote request:"+err.Error())
		return "", err
	}

	if resp.Body == nil {
		l.Logger.Log(ctx, "empty response from second")
		return "", nil
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	return string(out), err
}

func main() {

	r := chi.NewRouter()

	// Middleware добавляет GUID в context каждого HTTP запроса.
	//
	// Внутри:
	//   ctx = context.WithValue(ctx, guidKey, guid)
	//   req = req.WithContext(ctx)
	r.Use(tracker.Middleware)

	controller := Controller{
		Logic: LogicImpl{
			// Копирует GUID из context в header исходящего HTTP запроса
			RequestDecorator: tracker.Request,

			// Логгер извлекает GUID из context для корреляции логов
			Logger: tracker.Logger{},

			// адрес второго сервиса
			Remote: "http://localhost:4000",
		},
	}

	// endpoint:
	// GET /first?data=hello
	r.Get("/first", controller.First)

	http.ListenAndServe(":3000", r)
}
