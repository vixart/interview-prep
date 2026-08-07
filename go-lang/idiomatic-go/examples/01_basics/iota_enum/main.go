// iota: перечисления и битовые флаги.
// Правило: iota — это номер строки внутри блока const (с нуля),
// а строка без выражения повторяет ПРЕДЫДУЩЕЕ выражение.
package main

import "fmt"

type MailCategory int

const (
	Uncategorized  MailCategory = iota // 0
	Personal                           // 1
	Spam                               // 2
	Social                             // 3
	Advertisements                     // 4
)

// Значения повторяют предыдущее выражение, а iota все равно считает строки.
const (
	Field1 = 0        // 0
	Field2 = 1 + iota // 2  (iota == 1)
	Field3 = 20       // 20
	Field4            // 20 (повтор выражения "20")
	Field5 = iota     // 4  (iota == 4)
)

// Битовые флаги: сдвиг на iota.
type BitField int

const (
	Read    BitField = 1 << iota // 1
	Write                        // 2
	Execute                      // 4
	Delete                       // 8
)

func main() {
	fmt.Println(Uncategorized, Personal, Spam, Social, Advertisements)
	fmt.Println(Field1, Field2, Field3, Field4, Field5)

	perm := Read | Write
	fmt.Println("perm:", perm, "can write:", perm&Write != 0, "can execute:", perm&Execute != 0)

	// Значения iota — не часть API: вставка новой константы в середину
	// сдвинет все последующие. Не сохраняй их в БД/протокол как числа.
}
