// То же, что full_slice_expression, но БЕЗ ограничения емкости — и поведение
// становится непредсказуемым: append в y затирает элементы x, append в x и z
// дерутся за одну и ту же ячейку. Сравни вывод двух примеров.
package main

import "fmt"

func main() {
	x := make([]string, 0, 5)
	x = append(x, "a", "b", "c", "d")
	y := x[:2]
	// cap(y) == 5: y видит весь массив, а не только свои два элемента
	z := x[2:]
	fmt.Println(cap(x), cap(y), cap(z))
	y = append(y, "i", "j", "k")
	// пишем в чужую память: x[2], x[3] затерты, а x[4] заполнен
	x = append(x, "x")
	z = append(z, "y")
	// z стартует с x[2] и тоже пишет в общий массив — итог зависит от порядка строк
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)
}
