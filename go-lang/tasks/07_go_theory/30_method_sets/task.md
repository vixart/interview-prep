# 4.5.3. Исправление бага: реализация интерфейса и получатели методов

Тип задачи: ревью и исправление кода.

## Условие

> Исправьте код так, чтобы:
> 1. Реализации интерфейса работали корректно.
> 2. Программа завершалась успешно и без паники.
> 3. Лишние вызовы методов были устранены.
> 4. Вы смогли объяснить, где были ошибки и почему они возникли.

```go
package main

import (
    "errors"
    "fmt"
)

// PaymentProcessor - интерфейс для обработки платежей
type PaymentProcessor interface {
    Process(amount float64) error
    Verify(amount float64) bool
}

// CreditCardProcessor - реализация интерфейса для кредитной карты
type CreditCardProcessor struct {
    limit float64
}

// Process обрабатывает платеж
func (c *CreditCardProcessor) Process(amount float64) error {
    if amount > c.limit {
        return errors.New("...")
    }
    // ...
}

// Verify — объявлен с другим типом получателя (значение), из-за чего
// набор методов у CreditCardProcessor и *CreditCardProcessor различается
func (c CreditCardProcessor) Verify(amount float64) bool { /* ... */ }

// PayPalProcessor — вторая реализация с полем balance и аналогичной проблемой
type PayPalProcessor struct {
    balance float64
}

// ... Process / Verify для PayPalProcessor

// ExecutePayment вызывает методы Process и Verify
func ExecutePayment(processor PaymentProcessor, amount float64) {
    if processor.Verify(amount) {
        err := processor.Process(amount)
        if err != nil {
            fmt.Println("Error:", err)
        }
    } else {
        fmt.Println("Verification failed for amount:", amount)
    }
}

func main() {
    creditCard := CreditCardProcessor{limit: 100.0}
    payPal := PayPalProcessor{balance: 200.0}

    ExecutePayment(creditCard, 50.0)
    ExecutePayment(&creditCard, 50.0)
    ExecutePayment(&payPal, 150.0)
    ExecutePayment(payPal, 150.0)
}
```

## Ошибки компиляции, которые нужно объяснить

```
cannot use creditCard (variable of type CreditCardProcessor) as PaymentProcessor value in argument to ExecutePayment
cannot use payPal (variable of type PayPalProcessor) as PaymentProcessor value in argument to ExecutePayment
```

**Задание:** объясни ошибки компиляции и исправь код.
