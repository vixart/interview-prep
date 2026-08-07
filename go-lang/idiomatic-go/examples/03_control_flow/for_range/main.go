// Четыре формы for-range: индекс+значение, только значение (_), только ключ map
// и главное — forRangeIsACopy: переменная v это КОПИЯ, изменение v не меняет исходный срез.
package main

import "fmt"

func main() {
	forRangeKeyValue()
	forRangeIgnoreKey()
	forRangeMapKey()
	forRangeIsACopy()
}

func forRangeKeyValue() {
	fmt.Println("for-range loop, print index and value")
	evenVals := []int{2, 4, 6, 8, 10, 12}
	for i, v := range evenVals {
		// две переменные: индекс и КОПИЯ значения
		fmt.Println(i, v)
	}
}

func forRangeIgnoreKey() {
	fmt.Println("for-range loop, print value only")
	evenVals := []int{2, 4, 6, 8, 10, 12}
	for _, v := range evenVals {
		// индекс не нужен — гасим его через _
		fmt.Println(v)
	}
}

func forRangeMapKey() {
	fmt.Println("for-range loop over map, print key only")
	uniqueNames := map[string]bool{"Fred": true, "Raul": true, "Wilma": true}
	for k := range uniqueNames {
		// одна переменная при обходе map — это КЛЮЧ
		fmt.Println(k)
	}
}

func forRangeIsACopy() {
	fmt.Println("for-range loop, show that modifying value variable doesn't modify the original slice")
	evenVals := []int{2, 4, 6, 8, 10, 12}
	for _, v := range evenVals {
		v *= 2
		// меняем копию → исходный срез не изменится
	}
	fmt.Println(evenVals)
}
