// Как сделать свой долгий алгоритм прерываемым: цикл вычисления числа Пи
// на каждой итерации проверяет context.Cause(ctx) и выходит с ошибкой.
// Проверка контекста — обязанность самой функции, никто не прервет ее снаружи.
package main

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

func main() {
	calcWithTimeout(1)
	calcWithTimeout(2)
	calcWithTimeout(5)
	calcWithTimeout(10)
}

func calcWithTimeout(numSeconds time.Duration) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), numSeconds*time.Second)
	// разный лимит времени — разное количество посчитанных итераций
	defer cancelFunc()
	start := time.Now()
	result, err := calcPi(ctx)
	calcTime := time.Since(start)
	fmt.Println(result)
	fmt.Println(calcTime)
	fmt.Println(err)
}

func calcPi(ctx context.Context) (string, error) {
	var sum big.Float
	sum.SetInt64(0)
	var d big.Float
	d.SetInt64(1)
	two := big.NewFloat(2)
	i := 0
	for {
		if err := context.Cause(ctx); err != nil {
			// проверка контекста внутри своего цикла — иначе прервать вычисление невозможно
			fmt.Println("cancelled after", i, "iterations")
			return sum.Text('g', 100), err
		}
		var diff big.Float
		diff.SetInt64(4)
		diff.Quo(&diff, &d)
		if i%2 == 0 {
			sum.Add(&sum, &diff)
		} else {
			sum.Sub(&sum, &diff)
		}
		d.Add(&d, two)
		i++
	}
}
