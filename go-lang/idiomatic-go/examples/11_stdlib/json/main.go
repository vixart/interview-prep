// JSON туда и обратно потоково: json.NewEncoder пишет прямо в файл (io.Writer),
// json.NewDecoder читает прямо из файла (io.Reader) — без промежуточного []byte.
// Теги `json:"name"` задают имена полей.
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	err := ProcessPerson()
	if err != nil {
		slog.Error("error in processPerson", "msg", err)
	}
}

type Person struct {
	Name string `json:"name"`
	// тег задает имя поля в JSON; без тега берется имя поля как есть
	Age int `json:"age"`
}

func ProcessPerson() error {
	toFile := Person{
		Name: "Fred",
		Age:  40,
	}

	// Write it out
	tmpFile, err := os.CreateTemp(os.TempDir(), "sample-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	err = json.NewEncoder(tmpFile).Encode(toFile)
	// пишем прямо в io.Writer, без промежуточного []byte
	if err != nil {
		return err
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}

	// Read it back in again
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	var fromFile Person
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	// читаем прямо из io.Reader; передавать надо УКАЗАТЕЛЬ
	if err != nil {
		return err
	}
	err = tmpFile2.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", fromFile)
	return nil
}
