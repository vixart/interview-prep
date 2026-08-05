package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// countLetters читает данные из любого io.Reader и считает латинские буквы
func countLetters(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048) // буфер для чтения
	out := map[string]int{}   // результат: буква -> количество

	for {
		n, err := r.Read(buf) // читаем очередную порцию данных

		// обрабатываем только реально прочитанные байты
		for _, b := range buf[:n] {
			// проверяем что это латинская буква
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}

		// EOF означает, что поток закончился
		// важно: данные buf[:n] уже обработаны выше
		if err == io.EOF {
			return out, nil
		}

		// любая другая ошибка чтения
		if err != nil {
			return nil, err
		}
	}
}

// простой пример: читаем строку как поток
func simpleCountLetters() error {
	s := "The quick brown fox jumped over the lazy dog"

	// strings.NewReader превращает строку в io.Reader
	sr := strings.NewReader(s)

	counts, err := countLetters(sr)
	if err != nil {
		return err
	}

	fmt.Println(counts)
	return nil
}

// открывает gzip файл и возвращает:
// 1. gzip.Reader (для чтения распакованных данных)
// 2. функцию для закрытия ресурсов
func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
	r, err := os.Open(fileName) // открываем файл
	if err != nil {
		return nil, nil, err
	}

	gr, err := gzip.NewReader(r) // создаём reader для распаковки gzip
	if err != nil {
		return nil, nil, err
	}

	// возвращаем reader и функцию закрытия обоих ресурсов
	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}

// пример чтения из gzip файла
func gzipCountLetters() error {
	r, closer, err := buildGZipReader("my_data.txt.gz")
	if err != nil {
		return err
	}

	// гарантируем закрытие ресурсов
	defer closer()

	counts, err := countLetters(r)
	if err != nil {
		return err
	}

	fmt.Println(counts)
	return nil
}

func main() {
	// пример со строкой
	err := simpleCountLetters()
	if err != nil {
		slog.Error("error with simpleCountLetters", "msg", err)
	}

	// пример с gzip файлом
	err = gzipCountLetters()
	if err != nil {
		slog.Error("error with gzipCountLetters", "msg", err)
	}
}
