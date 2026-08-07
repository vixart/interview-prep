// Сквозной идентификатор запроса (GUID) через контекст: middleware кладет его в контекст,
// логгер достает и подмешивает в каждую строку лога, RequestDecorator прокидывает
// его в заголовке исходящего запроса к соседнему сервису. Это и есть законное
// применение context.WithValue — сквозные данные, а не аргументы функций.
package tracker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type guidKey int

const key guidKey = 1

// ключ неэкспортируемый — значение из контекста доступно только через функции пакета

func contextWithGUID(ctx context.Context, guid string) context.Context {
	return context.WithValue(ctx, key, guid)
}

func guidFromContext(ctx context.Context) (string, bool) {
	g, ok := ctx.Value(key).(string)
	return g, ok
}

func Middleware(h http.Handler) http.Handler {
	// middleware кладет GUID в контекст один раз для всего запроса
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if guid := req.Header.Get("X-GUID"); guid != "" {
			ctx = contextWithGUID(ctx, guid)
		} else {
			ctx = contextWithGUID(ctx, uuid.New().String())
			// нет входящего GUID — генерируем свой
		}
		req = req.WithContext(ctx)
		h.ServeHTTP(rw, req)
	})
}

type Logger struct{}

func (Logger) Log(ctx context.Context, message string) {
	if guid, ok := guidFromContext(ctx); ok {
		message = fmt.Sprintf("GUID: %s - %s", guid, message)
	}
	// do logging
	fmt.Println(message)
}

func Request(req *http.Request) *http.Request {
	ctx := req.Context()
	if guid, ok := guidFromContext(ctx); ok {
		req.Header.Add("X-GUID", guid)
	}
	return req
}
