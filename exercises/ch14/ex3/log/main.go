package log

import (
	"context"
	"fmt"
	"net/http"
)

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
)

// приватный тип для ключа — best practice,
// чтобы избежать коллизий ключей между пакетами
type logKey int

const (
	_ logKey = iota
	key
)

// кладём значение в context
// важно: context — immutable, возвращается НОВЫЙ ctx,
// старый при этом не изменяется
func ContextWithLevel(ctx context.Context, level Level) context.Context {
	return context.WithValue(ctx, key, level)
}

// достаём значение из context
// Value ищет по цепочке родителей (linked list),
// поэтому значение может быть установлено выше по стеку
func LevelFromContext(ctx context.Context) (Level, bool) {
	level, ok := ctx.Value(key).(Level)
	return level, ok
}

// Log читает уровень из context и принимает решение логировать или нет
// сам context здесь используется как "канал передачи метаданных"
func Log(ctx context.Context, level Level, message string) {
	var inLevel Level

	inLevel, ok := LevelFromContext(ctx)
	if !ok {
		// если в ctx ничего не положили — логирование отключено
		return
	}

	// context влияет на поведение функции (логируем или нет)
	if level == Debug && inLevel == Debug {
		fmt.Println(message)
	}

	if level == Info && (inLevel == Debug || inLevel == Info) {
		fmt.Println(message)
	}
}

// Middleware — точка входа, где мы "инжектим" данные в context запроса
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// достаём данные из внешнего мира (HTTP)
		level := r.URL.Query().Get("log_level")

		// создаём новый context на основе текущего request context
		// важно: сохраняем цепочку (timeout, cancel и т.д.)
		ctx := ContextWithLevel(r.Context(), Level(level))

		// подменяем context в request
		// дальше ВСЕ обработчики ниже будут видеть это значение
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
