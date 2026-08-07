// Идиома: функция, выделившая ресурс, возвращает closure для его освобождения,
// а вызывающий сразу ставит `defer closer()`. Так владение ресурсом остается явным.
// Запуск: go run . <путь к файлу>
package main

import (
	"io"
	"log"
	"os"
)

func getFile(name string) (*os.File, func(), error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	return file, func() {
		// вместе с ресурсом возвращаем функцию его освобождения
		file.Close()
	}, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no file specified")
	}
	f, closer, err := getFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer closer()
	// вызывающему остается один defer — он не обязан знать, что именно закрывается
	data := make([]byte, 2048)
	for {
		count, err := f.Read(data)
		os.Stdout.Write(data[:count])
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}
}
