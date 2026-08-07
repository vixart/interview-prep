// Тестируемая функция для табличного теста.
// Ветка "*" исправлена (в книге там был намеренный баг num1 + num2) — и заодно
// это урок: тест с данными 2 и 2 такой баг не поймал бы, потому что 2+2 == 2*2.
// Подбирай данные так, чтобы разные операции давали разные результаты.
package table

import (
	"errors"
	"fmt"
)

func DoMath(num1, num2 int, op string) (int, error) {
	switch op {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
		// в книге тут был баг (num1 + num2), который тест с данными 2 и 2 не ловит
	case "/":
		if num2 == 0 {
			return 0, errors.New("division by zero")
			// ошибочный случай тоже должен попасть в таблицу тестов
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("unknown operator %s", op)
	}
}
