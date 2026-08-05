package main

import (
	"bytes"         // буфер в памяти, реализует io.Writer и io.Reader
	"encoding/json" // пакет для работы с JSON
	"errors"        // для проверки ошибок типа io.EOF
	"fmt"           // вывод на экран
	"io"            // интерфейсы Reader/Writer и EOF
	"strings"       // создание io.Reader из строки
)

func main() {
	// JSON данные — несколько объектов подряд
	const data = `
		{"name": "Fred", "age": 40}
		{"name": "Mary", "age": 21}
		{"name": "Pat", "age": 30}
	`

	// структура для демаршалинга каждого JSON объекта
	var t struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// создаем JSON декодер, читающий из строки
	dec := json.NewDecoder(strings.NewReader(data)) // strings.NewReader → io.Reader

	// создаем буфер в памяти для сериализации обратно
	var b bytes.Buffer
	enc := json.NewEncoder(&b) // bytes.Buffer реализует io.Writer

	// цикл чтения объектов из JSON
	for {
		err := dec.Decode(&t) // читаем следующий JSON объект и заполняем t
		if err != nil {
			if errors.Is(err, io.EOF) { // EOF → больше данных нет
				break
			}
			panic(err) // любая другая ошибка → останов
		}

		fmt.Println(t)      // вывод прочитанного объекта
		err = enc.Encode(t) // сериализация обратно в JSON и запись в буфер
		if err != nil {
			panic(err)
		}
	}

	// получаем всю сериализованную строку из буфера
	out := b.String()
	fmt.Println(out) // выводим JSON всех объектов подряд
}
