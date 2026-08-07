// Множество на map[int]bool: дубликаты схлопываются (len(vals) != len(intSet)),
// проверка наличия — просто чтение (отсутствующий ключ дает false).
// Экономнее по памяти вариант — map[int]struct{} с идиомой «запятая-ok».
package main

import "fmt"

func main() {
	intSet := map[int]bool{}
	vals := []int{5, 10, 2, 5, 8, 7, 3, 9, 1, 2, 10}
	for _, v := range vals {
		intSet[v] = true
		// map как множество: повторные значения просто перезаписывают ключ
	}
	fmt.Println(len(vals), len(intSet))
	// 11 против 9 — дубликаты схлопнулись
	fmt.Println(intSet[5])
	fmt.Println(intSet[500])
	// отсутствующий ключ дает false — проверка наличия бесплатна
	if intSet[100] {
		fmt.Println("100 is in the set")
	}
}
