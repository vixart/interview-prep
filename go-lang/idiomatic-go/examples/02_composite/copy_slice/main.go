// copy(dst, src) копирует min(len(dst), len(src)) элементов и возвращает их количество.
// Важно: копируется по ДЛИНЕ приемника, емкость не учитывается.
package main

import "fmt"

func main() {
	x := []int{1, 2, 3, 4}
	y := make([]int, 4)
	num := copy(y, x)
	// копируется min(len(y), len(x)) = 4 элемента; num — сколько скопировано
	fmt.Println(y, num)
}
