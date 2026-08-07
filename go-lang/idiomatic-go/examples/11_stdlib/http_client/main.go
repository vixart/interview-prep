// HTTP-клиент как надо: свой http.Client с Timeout (у DefaultClient его нет),
// запрос через NewRequestWithContext, проверка StatusCode, defer res.Body.Close()
// и декодирование JSON прямо из тела ответа.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
		// у http.DefaultClient таймаута НЕТ — свой клиент нужен именно ради этого
	}

	req, err := http.NewRequestWithContext(context.Background(),
		// контекст в запросе = возможность отменить его снаружи
		http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-My-Client", "Learning Go")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	// тело ответа обязательно закрывать, иначе течет соединение
	if res.StatusCode != http.StatusOK {
		// ошибка от Do — только транспорт; HTTP-код проверяем сами
		panic(fmt.Sprintf("unexpected status: got %v", res.Status))
	}
	fmt.Println(res.Header.Get("Content-Type"))
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	// декодируем прямо из тела, не читая его целиком в память
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data)
}
