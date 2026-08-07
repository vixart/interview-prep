// Редкий оправданный goto: выход из цикла сразу к общему блоку завершения,
// минуя код, который выполняется только при нормальном завершении цикла.
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	a := rand.Intn(10)
	for a < 100 {
		if a%5 == 0 {
			goto done
			// выходим из цикла сразу к общему блоку завершения
		}
		a = a*2 + 1
	}
	fmt.Println("do something when the loop completes normally")
	// этот код выполнится только при НОРМАЛЬНОМ выходе из цикла — его goto и перепрыгивает
done:
	fmt.Println("do complicated stuff no matter why we left the loop")
	fmt.Println(a)
}
