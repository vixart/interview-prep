// Свой метод Is: тип с полем-срезом несравним через ==, поэтому errors.Is
// по умолчанию не нашел бы совпадение. Метод Is(target error) bool задает свое правило.
package main

import (
	"errors"
	"fmt"
	"slices"
)

type MyErr struct {
	Codes []int
	// поле-срез делает структуру НЕсравнимой → == не сработает
}

func (me MyErr) Error() string {
	return fmt.Sprintf("codes: %v", me.Codes)
}

func (me MyErr) Is(target error) bool {
	// свой метод Is задает правило сравнения; errors.Is вызовет именно его
	if me2, ok := target.(MyErr); ok {
		return slices.Equal(me.Codes, me2.Codes)
	}
	return false
}

func returnsWrappedMyErr() error {
	return fmt.Errorf("wrapping a MyErr: %w", MyErr{
		Codes: []int{2, 7, 1, 8, 2, 8},
	})
}

func main() {
	err := returnsWrappedMyErr()
	me := MyErr{Codes: []int{2, 7, 1, 8, 2, 8}}
	if errors.Is(err, me) {
		fmt.Println("found it!")
	}
}
