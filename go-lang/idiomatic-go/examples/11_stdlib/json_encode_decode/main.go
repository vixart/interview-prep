// Поток JSON-объектов: один Decoder читает подряд идущие объекты в цикле
// (конец потока — io.EOF), один Encoder пишет их в буфер.
// Так обрабатывают JSON Lines и длинные потоки, не загружая все в память.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func main() {
	const data = `
		{"name": "Fred", "age": 40}
		{"name": "Mary", "age": 21}
		{"name": "Pat", "age": 30}
	`
	var t struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	dec := json.NewDecoder(strings.NewReader(data))
	// один Decoder читает подряд идущие объекты из потока
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	for {
		err := dec.Decode(&t)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// io.EOF — нормальный конец потока, а не ошибка
				break
			}
			panic(err)
		}
		fmt.Println(t)
		err = enc.Encode(t)
		if err != nil {
			panic(err)
		}
	}
	out := b.String()
	fmt.Println(out)
}
