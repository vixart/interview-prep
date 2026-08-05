package main

import (
	"context"       // пакет для работы с контекстами (таймауты, отмена операций)
	"encoding/json" // пакет для работы с JSON
	"fmt"           // вывод в консоль
	"net/http"      // HTTP клиент и сервер
	"time"          // работа со временем
)

func main() {

	// Создаем HTTP-клиент.
	// Timeout ограничивает максимальное время ожидания ответа.
	// Если сервер зависнет — запрос автоматически прервется через 30 секунд.
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Создаем HTTP-запрос.
	// NewRequestWithContext принимает:
	// 1. context (для отмены запроса или таймаута)
	// 2. HTTP-метод
	// 3. URL
	// 4. тело запроса (io.Reader). Для GET оно не требуется → nil.
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://jsonplaceholder.typicode.com/todos/1",
		nil,
	)
	if err != nil {
		panic(err)
	}

	// Добавляем пользовательский HTTP-заголовок.
	// Заголовки часто используют для:
	// - авторизации
	// - версии API
	// - идентификации клиента
	req.Header.Add("X-My-Client", "Learning Go")

	// Отправляем HTTP-запрос.
	// Метод Do выполняет сетевой запрос и возвращает *http.Response.
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Тело ответа — поток (io.ReadCloser).
	// Его обязательно нужно закрывать после чтения,
	// иначе будет утечка соединений.
	defer res.Body.Close()

	// Проверяем код ответа сервера.
	// Обычно успешный ответ REST API — 200 OK.
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status: got %v", res.Status))
	}

	// Читаем один из HTTP-заголовков ответа.
	// Например Content-Type (тип данных, возвращенных сервером).
	fmt.Println(res.Header.Get("Content-Type"))

	// Анонимная структура для разбора JSON.
	// Теги json показывают, как поля JSON сопоставляются с полями Go.
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	// Создаем JSON-декодер, который читает напрямую из res.Body.
	// Decode автоматически:
	// - читает JSON
	// - парсит его
	// - заполняет структуру data
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	// Печатаем структуру с именами полей.
	// %+v выводит и имена полей, и их значения.
	fmt.Printf("%+v\n", data)
}
