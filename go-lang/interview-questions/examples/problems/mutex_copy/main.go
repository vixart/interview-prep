// Что происходит при копировании мьютекса.
//
// Короткий ответ: копируется его СОСТОЯНИЕ (в том числе «залочен»), и копия
// становится отдельной блокировкой. Две горутины начинают защищать разные
// мьютексы — взаимного исключения больше нет. Если скопировать залоченный
// мьютекс, копия навсегда останется залоченной → deadlock при Lock.
//
// Правило: sync.Mutex, RWMutex, WaitGroup, Once, atomic.* и все структуры,
// которые их содержат, передаются ТОЛЬКО по указателю.
//
// Плохой вариант вынесен в файл badcopy.go за тегом сборки, потому что
// `go vet` (анализатор copylocks) законно ругается на такой код:
//
//	go vet -tags copylocks ./problems/mutex_copy
//	go run -tags copylocks -race ./problems/mutex_copy
package main

import (
	"fmt"
	"sync"
)

// Counter содержит мьютекс, поэтому все методы — с УКАЗАТЕЛЬНЫМ приемником.
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

func main() {
	c := &Counter{} // работаем через указатель — копий не возникает
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	fmt.Println("правильно, через указатель:", c.Value()) // ровно 100

	runBadCopy() // в обычной сборке — заглушка, с тегом copylocks — демонстрация

	// Как это ловится:
	// - go vet: "passes lock by value" / "call of f copies lock value";
	// - go vet входит в go test, поэтому обычно ловится само;
	// - Mutex — это struct{state int32; sema uint32}, так что копия компилируется
	//   молча: без vet ошибку видно только по неверному поведению.
}
