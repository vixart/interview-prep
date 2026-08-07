// Буфер создается ОДИН раз и переиспользуется на каждой итерации чтения —
// идиома Go для чтения потоков без лишних аллокаций и работы сборщика мусора.
// Запуск: go run . <путь к файлу>
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	err := loadAndProcess(os.Args[1], func(data []byte) {
		fmt.Print(string(data))
	})
	if err != nil {
		fmt.Println("error:", err)
	}
}

func loadAndProcess(fileName string, process func([]byte)) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make([]byte, 100)
	// буфер выделяется ОДИН раз до цикла
	for {
		count, err := file.Read(data)
		// Read пишет в наш буфер и возвращает, сколько байт реально прочитал
		process(data[:count])
		// обрабатываем ровно прочитанную часть — и до проверки ошибки (данные приходят вместе с io.EOF)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if count == 0 {
			return nil
		}
	}
}
