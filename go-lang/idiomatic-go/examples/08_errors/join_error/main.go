// errors.Join собирает несколько независимых ошибок в одну —
// типичный случай валидации, где надо показать сразу все проблемы.
package main

import (
	"errors"
	"fmt"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func ValidatePerson(p Person) error {
	var errs []error
	if len(p.FirstName) == 0 {
		errs = append(errs, errors.New("field FirstName cannot be empty"))
		// копим все проблемы, а не выходим на первой
	}
	if len(p.LastName) == 0 {
		errs = append(errs, errors.New("field LastName cannot be empty"))
	}
	if p.Age < 0 {
		errs = append(errs, errors.New("field Age cannot be negative"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
		// Join склеивает несколько ошибок в одну; errors.Is найдет любую из них
	}
	return nil
}

func main() {
	err := ValidatePerson(Person{
		FirstName: "",
		LastName:  "",
		Age:       -1,
	})
	fmt.Println(err)
}
