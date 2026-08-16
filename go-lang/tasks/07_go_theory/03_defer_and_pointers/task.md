# 4.1.3. Defer и указатели: когда значения меняются неожиданно

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

Дан код. Нужно сказать, что выведет программа и в каком порядке.

```go
package main

import "fmt"

type Account struct {
    Balance int
}

func main() {
    initialBalance := 1000
    account := &Account{Balance: initialBalance}

    defer printBalance("Изначальный баланс", account.Balance)
    defer printBalance("Текущий баланс", account.Balance)
    defer printAccountBalance("Указатель на баланс", account)

    account.Balance += 500        // Делаем депозит
    updateBalance(account, 200)   // Снимаем средства
    account = &Account{Balance: 300} // Переназначаем указатель на новый аккаунт
}

func updateBalance(acc *Account, amount int) {
    acc.Balance -= amount
}

func printBalance(label string, balance int) {
    fmt.Printf("%s: %d\n", label, balance)
}

func printAccountBalance(label string, acc *Account) {
    fmt.Printf("%s: %d\n", label, acc.Balance)
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
