// map и срез передаются по значению, но значение — это дескриптор:
// изменения ЭЛЕМЕНТОВ видны вызывающему, а append внутри функции — нет
// (у копии среза своя длина).
package main

import "fmt"

func modMap(m map[int]string) {
	m[2] = "hello"
	// map передана по значению, но значение — дескриптор: правка видна снаружи
	m[3] = "goodbye"
	delete(m, 1)
}

func modSlice(s []int) {
	for k, v := range s {
		s[k] = v * 2
		// элементы среза меняются в той же памяти → снаружи видно
	}
	s = append(s, 10)
	// а append меняет ДЛИНУ у локальной копии дескриптора → снаружи не видно
}

func main() {
	m := map[int]string{
		1: "first",
		2: "second",
	}
	modMap(m)
	fmt.Println(m)

	s := []int{1, 2, 3}
	modSlice(s)
	fmt.Println(s)
}
