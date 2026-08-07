// Full slice expression x[low:high:max] ограничивает емкость подсреза.
// Из-за этого append в y и z выделяет новую память и НЕ затирает элементы x —
// это лечение проблемы из slice_append_storage.
package main

import "fmt"

func main() {
	x := make([]string, 0, 5)
	x = append(x, "a", "b", "c", "d")
	y := x[:2:2]
	// третье число обрезает cap: у y cap == 2, свободного места нет
	z := x[2:4:4]
	// у z тоже cap == 2
	fmt.Println(cap(x), cap(y), cap(z))
	y = append(y, "i", "j", "k")
	// cap исчерпан → append выделяет НОВЫЙ массив, x не пострадал
	x = append(x, "x")
	z = append(z, "y")
	// то же для z: без :max эта запись затерла бы элемент x
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)
}
