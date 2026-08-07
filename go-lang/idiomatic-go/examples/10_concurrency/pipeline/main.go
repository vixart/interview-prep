// Конвейер: этап AB считает две части параллельно, его результат уходит на этап C.
// Каждый этап умеет прерваться по контексту (общий таймаут 50 мс на весь конвейер).
// Запуск: go run . abc def
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type Input struct {
	A string
	B string
}

type COut struct {
	frequencyCount map[rune]int
}

func GatherAndProcess(ctx context.Context, data Input) (COut, error) {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	// общий дедлайн на ВЕСЬ конвейер, а не на отдельный этап
	defer cancel()

	ab := newABProcessor()
	ab.start(ctx, data)
	inputC, err := ab.wait(ctx)
	// ждем оба результата этапа AB (или ошибку/отмену)
	if err != nil {
		return COut{}, err
	}

	c := newCProcessor()
	c.start(ctx, inputC)
	// результат первого этапа — вход второго
	out, err := c.wait(ctx)
	return out, err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("expected input to be processed")
		os.Exit(1)
	}
	cout, err := GatherAndProcess(context.Background(), Input{
		A: os.Args[1],
		B: os.Args[2],
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cout)
}
