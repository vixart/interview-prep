package main

import (
	"encoding/json" // пакет для работы с JSON (Marshal / Unmarshal / Encoder / Decoder)
	"fmt"           // пакет для вывода
	"log/slog"      // пакет для логирования ошибок
	"os"            // работа с файловой системой
)

func main() {
	// вызываем функцию, которая пишет и читает JSON
	err := ProcessPerson()
	if err != nil {
		// логируем ошибку, если есть
		slog.Error("error in processPerson", "msg", err)
	}
}

// структура для демонстрации JSON
type Person struct {
	Name string `json:"name"` // имя поля в JSON — "name"
	Age  int    `json:"age"`  // имя поля в JSON — "age"
}

func ProcessPerson() error {
	// создаем пример данных для записи
	toFile := Person{
		Name: "Fred",
		Age:  40,
	}

	// -----------------------------
	// WRITE JSON TO FILE
	// -----------------------------

	// создаем временный файл
	tmpFile, err := os.CreateTemp(os.TempDir(), "sample-")
	if err != nil {
		return err
	}
	// гарантируем удаление файла при выходе из функции
	defer os.Remove(tmpFile.Name())

	// создаем JSON-энкодер, пишем данные структуры в файл
	err = json.NewEncoder(tmpFile).Encode(toFile)
	if err != nil {
		return err
	}

	// закрываем файл после записи
	err = tmpFile.Close()
	if err != nil {
		return err
	}

	// -----------------------------
	// READ JSON FROM FILE
	// -----------------------------

	// открываем файл для чтения
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}

	// структура для загрузки данных
	var fromFile Person

	// создаем JSON-декодер, читаем данные из файла в структуру
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	if err != nil {
		return err
	}

	// закрываем файл после чтения
	err = tmpFile2.Close()
	if err != nil {
		return err
	}

	// выводим прочитанные данные (Name и Age)
	fmt.Printf("%+v\n", fromFile)

	return nil
}
