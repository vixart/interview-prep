// MergeToMap: добавить в m[key] только новые значения из newValues.
//
// Чтобы не сравнивать каждый новый элемент со всеми старыми (O(n*m)),
// строим множество существующих значений — итог O(n + m).
// Мапа — ссылочный тип, поэтому функция меняет её на месте.
package main

import "fmt"

func MergeToMap(m map[string][]string, key string, newValues []string) {
	existing := make(map[string]struct{}, len(m[key]))
	for _, v := range m[key] {
		existing[v] = struct{}{}
	}

	for _, v := range newValues {
		if _, ok := existing[v]; ok {
			continue
		}
		m[key] = append(m[key], v)
		existing[v] = struct{}{} // защита от дубликатов внутри newValues
	}
}

func main() {
	m := map[string][]string{
		"group1": {"apple", "banana"},
		"group2": {"carrot"},
	}

	MergeToMap(m, "group1", []string{"banana", "cherry"})
	fmt.Println(m)
	// map[group1:[apple banana cherry] group2:[carrot]]

	MergeToMap(m, "group3", []string{"x", "x", "y"})
	fmt.Println(m)
	// map[group1:[apple banana cherry] group2:[carrot] group3:[x y]]
}
