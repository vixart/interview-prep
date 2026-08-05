package main

import (
	"encoding/json" // работа с JSON
	"fmt"           // вывод на экран
	"time"          // работа со временем
)

// -----------------------------
// Обычная структура Item для заказа
// -----------------------------
type Item struct {
	ID   string `json:"id"`   // имя поля в JSON — "id"
	Name string `json:"name"` // имя поля в JSON — "name"
}

// -----------------------------
// Структура Order
// -----------------------------
type Order struct {
	ID          string      `json:"id"`           // ID заказа
	Items       []Item      `json:"items"`        // список позиций
	DateOrdered RFC822ZTime `json:"date_ordered"` // дата заказа в специальном формате RFC822Z
	CustomerID  string      `json:"customer_id"`  // ID покупателя
}

// -----------------------------
// Новый тип времени RFC822ZTime
// -----------------------------
type RFC822ZTime struct {
	time.Time // встраиваем стандартное time.Time для доступа к методам
}

// -----------------------------
// Метод MarshalJSON для RFC822ZTime
// -----------------------------
// Позволяет сериализовать наш тип в JSON в формате RFC822Z
func (rt RFC822ZTime) MarshalJSON() ([]byte, error) {
	// форматируем внутреннее время по RFC822Z
	out := rt.Time.Format(time.RFC822Z)
	// оборачиваем в кавычки, так как JSON требует строки для даты
	return []byte(`"` + out + `"`), nil
}

// -----------------------------
// Метод UnmarshalJSON для RFC822ZTime
// -----------------------------
// Позволяет десериализовать JSON строку в наш тип времени
func (rt *RFC822ZTime) UnmarshalJSON(b []byte) error {
	// если в JSON значение null → оставляем поле пустым
	if string(b) == "null" {
		return nil
	}

	// разбираем строку времени по формату RFC822Z
	t, err := time.Parse(`"`+time.RFC822Z+`"`, string(b))
	if err != nil {
		return err
	}

	// сохраняем разобранное время в наш тип
	*rt = RFC822ZTime{t}
	return nil
}

// -----------------------------
// Основная функция
// -----------------------------
func main() {
	// JSON данные с датой в формате RFC822Z
	data := `
	{
		"id": "12345",
		"items": [
			{"id": "xyz123", "name": "Thing 1"},
			{"id": "abc789", "name": "Thing 2"}
		],
		"date_ordered": "01 May 20 13:01 +0000",
		"customer_id": "3"
	}`

	// создаем пустую структуру для загрузки данных
	var o Order

	// -----------------------------
	// Десериализация JSON в структуру
	// -----------------------------
	err := json.Unmarshal([]byte(data), &o)
	if err != nil {
		panic(err)
	}

	// вывод структуры, чтобы увидеть содержимое
	fmt.Printf("%+v\n", o)

	// можно использовать методы time.Time через RFC822ZTime
	fmt.Println(o.DateOrdered.Month()) // выводит месяц

	// -----------------------------
	// Сериализация обратно в JSON
	// -----------------------------
	out, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
