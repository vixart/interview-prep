package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/learning-go-book-2e/ch14/sample_code/context_user/identity"
)

// Logic — интерфейс бизнес-логики приложения.
//
// Важно: ctx передаётся первым аргументом.
// Это стандартная практика в Go.
//
// Контекст позволяет:
//   - отменять операции (если клиент закрыл соединение)
//   - передавать дедлайны
//   - передавать данные запроса (например user)
type Logic interface {
	BusinessLogic(ctx context.Context, user string, data string) (string, error)
}

// Controller — HTTP слой приложения.
// Он получает HTTP запросы и вызывает бизнес-логику.
type Controller struct {
	Logic Logic
}

// Login реализует максимально примитивную "аутентификацию".
//
// Пользователь передаётся через query:
//
// /login?user=Bob
//
// После этого имя пользователя записывается в cookie.
func (c Controller) Login(rw http.ResponseWriter, req *http.Request) {

	// Получаем имя пользователя из query параметра.
	userName := req.URL.Query().Get("user")

	// Проверяем что пользователь не пустой.
	if len(strings.TrimSpace(userName)) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("No user specified"))
		return
	}

	// identity.SetUser записывает cookie "identity".
	// Именно из этой cookie middleware позже извлечёт пользователя
	// и положит его в context.
	identity.SetUser(userName, rw)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("user logged in"))
}

// DoLogic — обработчик бизнес-запроса.
func (c Controller) DoLogic(rw http.ResponseWriter, req *http.Request) {

	// Каждый HTTP request в Go уже содержит context.
	//
	// Этот context создаётся сервером автоматически.
	// Он отменяется если:
	//   - клиент закрыл соединение
	//   - сработал timeout
	ctx := req.Context()

	// Достаём пользователя из context.
	//
	// Напомним цепочку:
	//
	// identity.Middleware ->
	//     context.WithValue(ctx, key, user)
	//
	// Поэтому здесь мы можем извлечь пользователя.
	user, ok := identity.UserFromContext(ctx)

	if !ok {
		// Это означает что middleware не добавил пользователя.
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получаем данные из запроса.
	data := req.URL.Query().Get("data")

	// Передаём context дальше в бизнес-логику.
	//
	// Это ОЧЕНЬ важный момент.
	//
	// context должен передаваться через все слои:
	//
	// HTTP handler -> service -> repository -> DB
	//
	// чтобы cancellation и deadlines работали корректно.
	result, err := c.Logic.BusinessLogic(ctx, user, data)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write([]byte(result))
}

// Logout — обработчик выхода пользователя.
func (c Controller) Logout(rw http.ResponseWriter, r *http.Request) {

	// Получаем context запроса.
	ctx := r.Context()

	// Проверяем что пользователь вообще есть в context.
	_, ok := identity.UserFromContext(ctx)

	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Удаляем cookie identity.
	identity.DeleteUser(rw)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("user logged out"))
}

// LogicImpl — конкретная реализация бизнес-логики.
type LogicImpl struct{}

// BusinessLogic получает context.
//
// В реальном приложении здесь могли бы быть:
//   - вызовы базы данных
//   - HTTP запросы к другим сервисам
//   - операции с дедлайнами
//
// context нужен чтобы:
//   - отменить операции если запрос отменён
//   - передать request-scoped данные
func (l LogicImpl) BusinessLogic(ctx context.Context, user string, data string) (string, error) {

	// В данном примере ctx не используется,
	// но его всё равно нужно принимать и передавать дальше
	// по архитектуре приложения.

	return fmt.Sprintf("Hello %s, thank you for sending me %s", user, data), nil
}

func main() {

	// Создаём новый роутер chi.
	r := chi.NewRouter()

	// Добавляем стандартное middleware логирования.
	r.Use(middleware.Logger)

	// Создаём контроллер.
	controller := Controller{
		Logic: LogicImpl{},
	}

	// endpoint:
	// GET /login?user=Bob
	r.Get("/login", controller.Login)

	// Группа маршрутов /business
	r.Route("/business", func(r chi.Router) {

		// r.With добавляет middleware только для этой группы.
		//
		// identity.Middleware делает следующее:
		//
		// 1. читает cookie identity
		// 2. извлекает user
		// 3. кладёт user в context
		//
		// ctx = context.WithValue(ctx, key, user)
		//
		// 4. создаёт новый request:
		//
		// req = req.WithContext(ctx)
		//
		// После этого все handler'ы внутри этой группы
		// получают request уже с user в context.
		r = r.With(identity.Middleware)

		// GET /business/?data=hello
		r.Get("/", controller.DoLogic)

		// GET /business/logout
		r.Get("/logout", controller.Logout)
	})

	// Запускаем HTTP сервер.
	http.ListenAndServe(":3000", r)
}
