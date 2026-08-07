// Обертывание всех ошибок функции одним defer.
// Работает только с ИМЕНОВАННЫМ возвращаемым значением ошибки:
// отложенная функция видит и может изменить его после return.
package main

import (
	"errors"
	"fmt"
)

func doThing1(v int) (int, error) { return v * 2, nil }
func doThing2(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty string")
	}
	return s + "!", nil
}
func doThing3(a int, b string) (string, error) { return fmt.Sprint(a, b), nil }

// Было бы три одинаковых fmt.Errorf("in DoSomeThings: %w", err) — стало одно.
func DoSomeThings(val1 int, val2 string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("in DoSomeThings: %w", err)
		}
	}()

	val3, err := doThing1(val1)
	if err != nil {
		return "", err
	}
	val4, err := doThing2(val2)
	if err != nil {
		return "", err
	}
	return doThing3(val3, val4)
}

func main() {
	fmt.Println(DoSomeThings(5, "hello"))
	_, err := DoSomeThings(5, "")
	fmt.Println(err)                                // in DoSomeThings: empty string
	fmt.Println(errors.Unwrap(err))                 // empty string
	fmt.Println(errors.Is(err, errors.Unwrap(err))) // true
}
