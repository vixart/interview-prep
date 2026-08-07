// append в подсрез затирает элемент РОДИТЕЛЯ.
// y := x[:2] имеет len 2, но cap 4, поэтому append(y, "z") пишет в ту же память
// и x становится [a b z d]. Самая частая ловушка срезов.
package main

import "fmt"

func main() {
	x := []string{"a", "b", "c", "d"}
	y := x[:2]
	// len(y) == 2, но cap(y) == 4 — хвост родителя все еще доступен
	fmt.Println(cap(x), cap(y))
	y = append(y, "z")
	// места хватает → пишем в тот же массив и затираем x[2]
	fmt.Println("x:", x)
	// x стал [a b z d] — хотя x мы не трогали
	fmt.Println("y:", y)
}
