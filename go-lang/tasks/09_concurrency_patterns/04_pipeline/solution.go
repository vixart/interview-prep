// Pipeline: конвейер обработки транзакций.
//
// Каждый этап — горутина: читает входной канал, пишет в выходной и
// ЗАКРЫВАЕТ его за собой (правило: канал закрывает тот, кто пишет).
// Все этапы умеют выходить по ctx.Done() — иначе при раннем выходе
// потребителя горутины этапов зависнут навсегда.
package main

import (
	"context"
	"fmt"
)

type Transaction struct {
	ID       int
	Amount   float64
	Currency string
}

const usdRate = 90.0 // условный курс RUB -> USD

// 1. Чтение транзакций из исходных данных.
func generator(ctx context.Context, txs []Transaction) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		defer close(out)
		for _, tx := range txs {
			select {
			case out <- tx:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// 2. Фильтрация: убираем транзакции с отрицательными суммами.
func filterNegative(ctx context.Context, in <-chan Transaction) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		defer close(out)
		for tx := range in {
			if tx.Amount < 0 {
				continue
			}
			select {
			case out <- tx:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// 3. Конвертация валюты в доллары.
func toUSD(ctx context.Context, in <-chan Transaction) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		defer close(out)
		for tx := range in {
			if tx.Currency == "RUB" {
				tx.Amount /= usdRate
				tx.Currency = "USD"
			}
			select {
			case out <- tx:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	txs := []Transaction{
		{ID: 1, Amount: 9000, Currency: "RUB"},
		{ID: 2, Amount: -500, Currency: "RUB"}, // отфильтруется
		{ID: 3, Amount: 100, Currency: "USD"},
		{ID: 4, Amount: 4500, Currency: "RUB"},
	}

	// 4. Сохранение результатов: собираем обработанные транзакции.
	var result []Transaction
	for tx := range toUSD(ctx, filterNegative(ctx, generator(ctx, txs))) {
		result = append(result, tx)
	}

	for _, tx := range result {
		fmt.Printf("#%d %.2f %s\n", tx.ID, tx.Amount, tx.Currency)
	}
	// #1 100.00 USD
	// #3 100.00 USD
	// #4 50.00 USD
}
