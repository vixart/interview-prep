// Свой тип с несколькими ошибками внутри: Unwrap() []error (не error!).
// Внизу — шаблон разбора: switch по интерфейсам { Unwrap() error } и { Unwrap() []error },
// чтобы обработать и одиночную, и множественную обертку.
package main

import (
	"errors"
	"fmt"
)

type Status int

type StatusErr struct {
	Status  Status
	Message string
}

func (se StatusErr) Error() string {
	return se.Message
}

const (
	InvalidLogin Status = iota + 1
	NotFound
)

type MyError struct {
	Code   int
	Errors []error
}

func (m MyError) Error() string {
	return errors.Join(m.Errors...).Error()
}

func (m MyError) Unwrap() []error {
	// СРЕЗ ошибок, а не одна — вторая форма Unwrap
	return m.Errors
}

func funcThatReturnsAnError() error {
	return MyError{
		Code: 12,
		Errors: []error{
			StatusErr{
				Status:  NotFound,
				Message: "file Not Found",
			},
			errors.New("a simple string error"),
		},
	}
}

func main() {
	var err error
	err = funcThatReturnsAnError()
	switch err := err.(type) {
	case interface{ Unwrap() error }:
		// шаблон разбора: одиночная обертка...
		// handle single error
		innerErr := err.Unwrap()
		// process innerErr
		fmt.Println(innerErr)
	case interface{ Unwrap() []error }:
		// ...и множественная
		//handle multiple wrapped errors
		innerErrs := err.Unwrap()
		for _, innerErr := range innerErrs {
			// process each innerErr
			fmt.Println(innerErr)
		}
	default:
		// handle no wrapped error
	}
}
