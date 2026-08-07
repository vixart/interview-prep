// Сигнальная ошибка: пакет archive/zip экспортирует zip.ErrFormat,
// и вызывающий сверяет с ней результат. Так помечают ситуацию,
// в которой дальнейшая обработка невозможна. В новом коде вместо == пишут errors.Is.
package main

import (
	"archive/zip"
	"bytes"
	"fmt"
)

func main() {
	data := []byte("This is not a zip file")
	notAZipFile := bytes.NewReader(data)
	_, err := zip.NewReader(notAZipFile, int64(len(data)))
	if err == zip.ErrFormat {
		// сигнальная ошибка объявлена в пакете как переменная и сверяется по значению
		fmt.Println("Told you so")
	}
}
