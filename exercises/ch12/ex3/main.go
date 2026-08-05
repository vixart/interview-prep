package main

import (
	"fmt"
	"math"
	"sync"
)

func buildSquareRootMap() map[int]float64 {
	out := make(map[int]float64, 100_000)
	for i := 0; i < 100_000; i++ {
		out[i] = math.Sqrt(float64(i))
	}
	return out
}

var sqrtsCache = sync.OnceValue(buildSquareRootMap)

func main() {
	for i := 0; i < 100_000; i += 1_000 {
		fmt.Println(i, sqrtsCache()[i])
	}
}
