package tracker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// собственный тип ключа для context.
// Использование уникального типа предотвращает коллизии ключей
// с другими пакетами, которые тоже могут использовать context.
type guidKey int

const key guidKey = 1

// contextWithGUID добавляет GUID в context.
// context.WithValue не изменяет существующий context,
// а возвращает новый context, который ссылается на старый
// и хранит дополнительное значение.
func contextWithGUID(ctx context.Context, guid string) context.Context {
	return context.WithValue(ctx, key, guid)
}

// guidFromContext извлекает GUID из context.
// ctx.Value возвращает interface{}, поэтому нужен type assertion.
func guidFromContext(ctx context.Context) (string, bool) {
	g, ok := ctx.Value(key).(string)
	return g, ok
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// каждый HTTP request уже содержит context,
		// созданный HTTP сервером.
		ctx := req.Context()

		// если клиент передал GUID — используем его
		if guid := req.Header.Get("X-GUID"); guid != "" {
			ctx = contextWithGUID(ctx, guid)
		} else {
			// иначе генерируем новый GUID для трассировки запроса
			ctx = contextWithGUID(ctx, uuid.New().String())
		}

		// Request immutable → WithContext создаёт новый request
		req = req.WithContext(ctx)

		// передаём request дальше по цепочке handler'ов
		h.ServeHTTP(rw, req)
	})
}

type Logger struct{}

// Log получает context, чтобы извлечь GUID текущего запроса.
// Благодаря context логгер не зависит от HTTP и может
// использоваться в любом слое приложения.
func (Logger) Log(ctx context.Context, message string) {
	if guid, ok := guidFromContext(ctx); ok {
		message = fmt.Sprintf("GUID: %s - %s", guid, message)
	}

	// пример логирования
	fmt.Println(message)
}

// Request добавляет GUID из context в исходящий HTTP запрос.
// Это позволяет прокинуть идентификатор запроса в другой сервис
// (distributed tracing).
func Request(req *http.Request) *http.Request {
	ctx := req.Context()

	if guid, ok := guidFromContext(ctx); ok {
		req.Header.Add("X-GUID", guid)
	}

	return req
}
