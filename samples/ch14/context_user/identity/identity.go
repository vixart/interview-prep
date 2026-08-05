package identity

import (
	"context"
	"net/http"
)

// userKey — это тип ключа, который используется для хранения значения в context.
// В Go крайне важно использовать собственный тип ключа,
// чтобы избежать коллизий с ключами из других пакетов.
//
// Если бы мы использовали string ("user"), другой пакет мог бы
// случайно использовать тот же ключ и перезаписать значение.
type userKey int

const (
	// iota используется для генерации уникальных значений.
	// Первый элемент игнорируем (_), чтобы "key" не был равен 0.
	_ userKey = iota
	key
)

// ContextWithUser кладёт пользователя в context.
//
// context в Go — это неизменяемая структура.
// Поэтому context.WithValue НЕ модифицирует старый context,
// а создаёт новый context, который ссылается на старый.
//
// Структура получается примерно такая:
//
// parentCtx
//
//	│
//	▼
//
// ctxWithUser (value: key -> user)
//
// Поэтому мы должны вернуть новый context.
func ContextWithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, key, user)
}

// UserFromContext извлекает пользователя из context.
//
// ctx.Value(key) возвращает interface{}.
// Поэтому мы должны сделать type assertion -> (string).
//
// ok показывает:
//
//	true  — значение было и оно string
//	false — значения нет или тип другой
//
// Такой паттерн стандартный для извлечения данных из context.
func UserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(key).(string)
	return user, ok
}

// extractUser извлекает пользователя из HTTP запроса.
//
// В данном примере пользователь хранится в cookie "identity".
//
// В реальном приложении обычно:
//   - cookie подписывается
//   - используется JWT
//   - или session store
//
// чтобы пользователь не мог подделать значение.
func extractUser(req *http.Request) (string, error) {

	// req.Cookie ищет cookie по имени
	userCookie, err := req.Cookie("identity")
	if err != nil {
		return "", err
	}

	return userCookie.Value, nil
}

// Middleware — это HTTP middleware.
//
// Middleware — это функция, которая:
//
//	принимает handler
//	возвращает новый handler
//
// Этот новый handler выполняется ПЕРЕД оригинальным handler.
//
// Здесь middleware:
// 1. извлекает пользователя из cookie
// 2. кладёт его в context
// 3. передаёт управление следующему handler
func Middleware(h http.Handler) http.Handler {

	// http.HandlerFunc — адаптер, превращающий функцию
	// в объект http.Handler.
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// Пытаемся получить пользователя из cookie.
		user, err := extractUser(req)
		if err != nil {

			// Если cookie нет — пользователь не авторизован.
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("unauthorized"))
			return
		}

		// Получаем текущий context запроса.
		//
		// В Go каждый HTTP запрос уже имеет context.
		// Этот context:
		//   - отменяется если клиент закрыл соединение
		//   - используется для дедлайнов
		//   - передаётся через всю цепочку вызовов.
		ctx := req.Context()

		// Создаём новый context, в который кладём user.
		//
		// Теперь в context есть значение:
		//
		// key -> user
		ctx = ContextWithUser(ctx, user)

		// Request в Go тоже immutable.
		// Поэтому WithContext создаёт НОВЫЙ request.
		req = req.WithContext(ctx)

		// Передаём управление следующему handler.
		//
		// Теперь любой код ниже по цепочке может сделать:
		//
		// user, _ := identity.UserFromContext(req.Context())
		//
		// и получить текущего пользователя.
		h.ServeHTTP(rw, req)
	})
}

// SetUser устанавливает cookie с идентификатором пользователя.
//
// Обычно вызывается после успешного логина.
func SetUser(user string, rw http.ResponseWriter) {

	http.SetCookie(rw, &http.Cookie{
		Name:  "identity",
		Value: user,
	})
}

// DeleteUser удаляет cookie.
//
// Это делается через:
// MaxAge = -1
//
// Браузер сразу удаляет cookie.
func DeleteUser(rw http.ResponseWriter) {

	http.SetCookie(rw, &http.Cookie{
		Name:   "identity",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
