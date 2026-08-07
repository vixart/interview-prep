// Набор методов: у *Counter есть и Increment (указательный приемник), и String,
// у значения Counter — только String. Поэтому значение не реализует Incrementer
// (строка закомментирована), а указатель — реализует.
package main

import (
	"fmt"
	"time"
)

type Counter struct {
	total       int
	lastUpdated time.Time
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdated = time.Now()
}

func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdated)
}

type Incrementer interface {
	Increment()
}

func main() {
	var myStringer fmt.Stringer
	var myIncrementer Incrementer
	pointerCounter := &Counter{}
	valueCounter := Counter{}

	myStringer = pointerCounter // ok
	myStringer = valueCounter   // ok
	// значение реализует интерфейс: String объявлен с приемником-значением
	myIncrementer = pointerCounter // ok
	// у *Counter в наборе методов есть и Increment, и String
	// Набор методов значимого типа не включает методы с указательным приемником:
	// myIncrementer = valueCounter // ошибка компиляции

	fmt.Println(myStringer, myIncrementer)
}
