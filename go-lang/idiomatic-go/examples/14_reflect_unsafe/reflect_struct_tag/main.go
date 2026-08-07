// Основа рефлексии: reflect.TypeOf → NumField/Field(i) → имя поля, его тип и ТЕГИ.
// Ровно так encoding/json и ORM понимают, в какую колонку/ключ класть поле.
package main

import (
	"fmt"
	"reflect"
)

func main() {
	type Foo struct {
		A int `myTag:"value"`
		// тег — просто строка; смысл ей придает тот, кто ее читает
		B string `myTag:"value2"`
	}

	var f Foo
	ft := reflect.TypeOf(f)
	// reflect.Type — описание типа во время выполнения
	for i := 0; i < ft.NumField(); i++ {
		// NumField/Field работают ТОЛЬКО для Kind() == Struct, иначе паника
		curField := ft.Field(i)
		fmt.Println(curField.Name, curField.Type.Name(),
			curField.Tag.Get("myTag"))
		// ровно так encoding/json достает `json:"..."`
	}
}
