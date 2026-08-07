// Ловушка: функция объявила возвращаемый тип как СВОЙ тип ошибки (или вернула
// неинициализированную переменную этого типа) — и `err != nil` становится true,
// хотя ошибки не было. Причина: интерфейс хранит тип, поэтому не равен nil.
// Правильно — возвращать error и явный nil.
package main

import (
	"fmt"
)

type Status int

const (
	InvalidLogin Status = iota + 1
	NotFound
)

type StatusErr struct {
	Status  Status
	Message string
}

func (se StatusErr) Error() string {
	return se.Message
}

func GenerateErrorBroken(flag bool) error {
	var genErr StatusErr
	// ЛОВУШКА: пустая структура, но не nil
	if flag {
		genErr = StatusErr{
			Status: NotFound,
		}
	}
	return genErr
	// кладем ее в error → интерфейс получает тип, и err != nil даже без ошибки
}

func GenerateErrorOKReturnNil(flag bool) error {
	if flag {
		return StatusErr{
			Status: NotFound,
		}
	}
	return nil
	// правильно: явный nil
}

func GenerateErrorUseErrorVar(flag bool) error {
	var genErr error
	if flag {
		genErr = StatusErr{
			Status: NotFound,
		}
	}
	return genErr
}

func main() {
	err := GenerateErrorBroken(true)
	fmt.Println("GenerateErrorBroken(true) returns non-nil error:", err != nil)
	err = GenerateErrorBroken(false)
	fmt.Println("GenerateErrorBroken(false) returns non-nil error:", err != nil)
	err = GenerateErrorOKReturnNil(true)
	fmt.Println("GenerateErrorOKReturnNil(true) returns non-nil error:", err != nil)
	err = GenerateErrorOKReturnNil(false)
	fmt.Println("GenerateErrorOKReturnNil(false) returns non-nil error:", err != nil)
	err = GenerateErrorUseErrorVar(true)
	fmt.Println("GenerateErrorUseErrorVar(true) returns non-nil error:", err != nil)
	err = GenerateErrorUseErrorVar(false)
	fmt.Println("GenerateErrorUseErrorVar(false) returns non-nil error:", err != nil)
}
