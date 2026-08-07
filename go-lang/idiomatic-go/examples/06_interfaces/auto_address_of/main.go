// Метод с указательным приемником можно вызвать на обычной переменной:
// Go сам подставляет &c, потому что переменная адресуема.
// Из-за этого разница между приемниками не бросается в глаза — до передачи в интерфейс
// (см. method_set) или в функцию (см. value_vs_pointer_receiver).
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
	// указательный приемник — метод меняет структуру
	c.total++
	c.lastUpdated = time.Now()
}

func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdated)
}

func main() {
	var c Counter
	fmt.Println(c.String())
	c.Increment()
	// c — значение, но Go сам подставляет (&c).Increment(), потому что c адресуема
	fmt.Println(c.String())
}
