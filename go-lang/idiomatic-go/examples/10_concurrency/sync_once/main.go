// Ленивая инициализация: тяжелый парсер создается ровно один раз даже при
// конкурентных вызовах. sync.Once объявлен на уровне пакета — внутри функции
// он создавался бы заново на каждый вызов и смысла не имел.
package main

import (
	"fmt"
	"sync"
)

func main() {
	// "initializing!" will print out only once
	result := Parse("hello")
	fmt.Println(result)
	result2 := Parse("goodbye")
	fmt.Println(result2)
}

type SlowComplicatedParser interface {
	Parse(string) string
}

var parser SlowComplicatedParser
var once sync.Once

// на уровне пакета: внутри функции создавался бы новый Once на каждый вызов

func Parse(dataToParse string) string {
	once.Do(func() {
		// выполнится ровно один раз, даже при конкурентных вызовах Parse
		parser = initParser()
	})
	return parser.Parse(dataToParse)
}

func initParser() SlowComplicatedParser {
	// do all sorts of setup and loading here
	fmt.Println("initializing!")
	return SCPI{}
}

type SCPI struct {
}

func (s SCPI) Parse(in string) string {
	if len(in) > 1 {
		return in[0:1]
	}
	return ""
}
