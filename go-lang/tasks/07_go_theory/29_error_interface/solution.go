// Вернуть ошибку без импорта errors/fmt.
//
// error — это обычный интерфейс из встроенного пакета:
//
//	type error interface {
//	    Error() string
//	}
//
// Достаточно объявить свой тип с методом Error() string.
package main

type myError struct {
	msg string
}

func (e *myError) Error() string { return e.msg }

func handle() error {
	return &myError{msg: "что-то пошло не так"}
}

func main() {
	err := handle()
	// println печатает интерфейс как пару указателей (itab, data),
	// поэтому выводим текст явно:
	println(err.Error()) // что-то пошло не так

	// Классическая ловушка: типизированный nil-указатель в интерфейсе.
	// Если бы handle возвращала (*myError)(nil) как error,
	// проверка err != nil была бы ИСТИННОЙ — интерфейс не пуст,
	// в нём лежит пара (тип=*myError, значение=nil).
	var typedNil *myError
	var iface error = typedNil
	println(iface == nil) // false!
}
