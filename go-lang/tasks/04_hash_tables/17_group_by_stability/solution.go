// Группировка серверов по стабильности: map[стабильность] -> список ID.
//
// Ключ — строковое представление float64 без лишних нулей
// (strconv.FormatFloat с prec=-1), поэтому 97 и 97.1 — разные ключи.
// Время O(n).
package main

import (
	"fmt"
	"strconv"
)

type ServerStat struct {
	Server    int
	Stability float64
}

func groupByStability(stats []ServerStat) map[string][]int {
	result := make(map[string][]int)

	for _, s := range stats {
		key := strconv.FormatFloat(s.Stability, 'f', -1, 64)
		result[key] = append(result[key], s.Server)
	}

	return result
}

func main() {
	stats := []ServerStat{
		{Server: 1, Stability: 99},
		{Server: 2, Stability: 97},
		{Server: 3, Stability: 34},
		{Server: 4, Stability: 97},
		{Server: 5, Stability: 97.1},
	}

	fmt.Println(groupByStability(stats))
	// map[34:[3] 97:[2 4] 97.1:[5] 99:[1]]
}
