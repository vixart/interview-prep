// io.Reader на практике: countLetters работает с ЛЮБЫМ источником (строка, файл, gzip).
// Ключевые моменты: буфер выделяется один раз; сначала обрабатываются прочитанные buf[:n],
// и только потом проверяется io.EOF (последние байты приходят вместе с ним);
// gzip.NewReader — декоратор поверх другого Reader.
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func countLetters(r io.Reader) (map[string]int, error) {
	// работает с ЛЮБЫМ источником: строкой, файлом, сетью, gzip
	buf := make([]byte, 2048)
	// буфер один на весь цикл — Read пишет в него, а не выделяет новый
	out := map[string]int{}
	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] {
			// СНАЧАЛА обрабатываем прочитанные n байт...
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		if err == io.EOF {
			// ...и только потом смотрим на ошибку: последние байты приходят вместе с EOF
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func simpleCountLetters() error {
	s := "The quick brown fox jumped over the lazy dog"
	sr := strings.NewReader(s)
	// строка тоже io.Reader — функцию не пришлось менять
	counts, err := countLetters(sr)
	if err != nil {
		return err
	}
	fmt.Println(counts)
	return nil
}

func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	gr, err := gzip.NewReader(r)
	// декоратор: Reader, обернутый в Reader
	if err != nil {
		return nil, nil, err
	}
	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}

func gzipCountLetters() error {
	r, closer, err := buildGZipReader("my_data.txt.gz")
	if err != nil {
		return err
	}
	defer closer()
	counts, err := countLetters(r)
	if err != nil {
		return err
	}
	fmt.Println(counts)
	return nil
}

func main() {
	err := simpleCountLetters()
	if err != nil {
		slog.Error("error with simpleCountLetters", "msg", err)
	}

	err = gzipCountLetters()
	if err != nil {
		slog.Error("error with gzipCountLetters", "msg", err)
	}
}
