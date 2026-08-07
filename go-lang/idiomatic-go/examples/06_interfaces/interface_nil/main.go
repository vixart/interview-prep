// Типизированный nil — главная ловушка интерфейсов.
// Указатель nil, интерфейс nil, но после присваивания указателя в интерфейс
// он перестает быть nil: тип внутри интерфейса задан, значит пара (тип, значение) != (nil, nil).
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
	var pointerCounter *Counter
	// указатель nil
	fmt.Println(pointerCounter == nil) // prints true
	var incrementer Incrementer
	fmt.Println(incrementer == nil) // prints true
	incrementer = pointerCounter
	// кладем nil-указатель в интерфейс: тип = *Counter, значение = nil
	fmt.Println(incrementer == nil) // prints false
	// ВОТ ОНО: интерфейс не nil, потому что тип внутри задан
}
