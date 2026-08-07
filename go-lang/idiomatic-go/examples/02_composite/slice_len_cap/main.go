// Рост среза: как меняются len и cap при append, начиная с nil-среза.
// Видно удвоение емкости (0 → 1 → 2 → 4 → 8) и то, что append в nil-срез легален.
package main

import "fmt"

func main() {
	var x []int
	// nil-срез: len 0, cap 0, но append работать будет
	fmt.Println(x, len(x), cap(x))
	x = append(x, 10)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 20)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 30)
	// cap исчерпан → новый массив вдвое больше, данные скопированы
	fmt.Println(x, len(x), cap(x))
	x = append(x, 40)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 50)
	fmt.Println(x, len(x), cap(x))
}
