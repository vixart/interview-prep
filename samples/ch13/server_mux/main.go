package main

import (
	"net/http" // пакет для создания HTTP-серверов и клиентов
	"time"     // пакет для работы со временем (нужен для таймаутов)
)

/*
Endpoints (доступные URL):

/person/greet
/dog/greet
/hello

Пример запросов:

GET http://localhost:8080/person/greet
GET http://localhost:8080/dog/greet
GET http://localhost:8080/hello
*/

func main() {

	// -------------------------------
	// Router для /person/*
	// -------------------------------

	// Создаем отдельный ServeMux.
	// Он будет обрабатывать только маршруты, связанные с "person".
	person := http.NewServeMux()

	// Регистрируем обработчик для пути "/greet".
	// HandleFunc — удобная версия Handle, принимающая обычную функцию.
	person.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {

		// Отправляем HTTP-ответ клиенту.
		w.Write([]byte("greetings!\n"))
	})

	// -------------------------------
	// Router для /dog/*
	// -------------------------------

	// Еще один независимый маршрутизатор.
	dog := http.NewServeMux()

	// Регистрируем обработчик "/greet".
	dog.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {

		// Ответ клиенту
		w.Write([]byte("good puppy!\n"))
	})

	// -------------------------------
	// Главный router
	// -------------------------------

	// Главный ServeMux будет принимать все входящие запросы.
	mux := http.NewServeMux()

	// Регистрируем подмаршрут /person/*
	//
	// StripPrefix удаляет "/person" из URL перед тем,
	// как передать запрос во внутренний router "person".
	//
	// Например:
	//
	// клиент: /person/greet
	// StripPrefix → /greet
	// person mux получает: /greet
	mux.Handle("/person/", http.StripPrefix("/person", person))

	// Аналогично для /dog/*
	mux.Handle("/dog/", http.StripPrefix("/dog", dog))

	// -------------------------------
	// Простой endpoint
	// -------------------------------

	// Этот маршрут обрабатывается прямо главным mux.
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {

		// Отправляем простой текстовый ответ
		w.Write([]byte("Hello!\n"))
	})

	// -------------------------------
	// Конфигурация HTTP сервера
	// -------------------------------

	s := http.Server{

		// Сервер слушает порт 8080 на всех интерфейсах
		Addr: ":8080",

		// Максимальное время чтения запроса
		ReadTimeout: 30 * time.Second,

		// Максимальное время записи ответа
		WriteTimeout: 90 * time.Second,

		// Время ожидания следующего запроса при keep-alive
		IdleTimeout: 120 * time.Second,

		// Главный обработчик всех HTTP-запросов
		Handler: mux,
	}

	// Запускаем сервер:
	// 1. открывается TCP порт
	// 2. сервер начинает слушать запросы
	// 3. каждый запрос передается mux
	err := s.ListenAndServe()

	// Проверяем ошибку запуска/работы сервера
	if err != nil {

		// ErrServerClosed — нормальная ошибка при штатном завершении
		if err != http.ErrServerClosed {

			// Любая другая ошибка — критическая
			panic(err)
		}
	}
}
