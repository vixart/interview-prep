// Вариативный параметр: внутри функции это обычный срез, снаружи можно передать
// любое количество аргументов или развернуть готовый срез через `s...`.
package main

import "fmt"

func addTo(base int, vals ...int) []int {
	// вариативный параметр всегда последний; внутри это обычный []int
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}

func main() {
	fmt.Println(addTo(3))
	// можно вообще не передавать: vals будет пустым
	fmt.Println(addTo(3, 2))
	fmt.Println(addTo(3, 2, 4, 6, 8))
	a := []int{4, 3}
	fmt.Println(addTo(3, a...))
	// готовый срез разворачивается через ...
	fmt.Println(addTo(3, []int{1, 2, 3, 4, 5}...))
}
