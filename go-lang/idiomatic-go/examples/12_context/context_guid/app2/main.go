// Второй сервис цепочки: берет GUID из ЗАГОЛОВКА входящего запроса
// (его положил app1), кладет в свой контекст и логирует с ним же.
// Так один идентификатор связывает логи двух сервисов.
package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"interviewprep/examples/12_context/context_guid/tracker"
	"net/http"
)

type Logic interface {
	QueryHandler(ctx context.Context, query string) (string, error)
}
type Controller struct {
	Logic Logic
}

func (c Controller) Second(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	query := req.URL.Query().Get("query")
	result, err := c.Logic.QueryHandler(ctx, query)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write([]byte(result))
}

type Logger interface {
	Log(context.Context, string)
}

type LogicImpl struct {
	Logger Logger
	Remote string
}

func (l LogicImpl) QueryHandler(ctx context.Context, query string) (string, error) {
	l.Logger.Log(ctx, "starting QueryHandler with query: "+query)
	// логгер сам вытащит GUID из контекста — в лог попадет тот же идентификатор, что у app1
	return fmt.Sprintf("got query: '%s' from first", query), nil
}

func main() {
	r := chi.NewRouter()
	r.Use(tracker.Middleware)
	// middleware достанет GUID из заголовка X-GUID (его прислал app1) и положит в контекст
	controller := Controller{
		Logic: LogicImpl{
			Logger: tracker.Logger{},
		},
	}
	r.Get("/second", controller.Second)
	http.ListenAndServe(":4000", r)
}
