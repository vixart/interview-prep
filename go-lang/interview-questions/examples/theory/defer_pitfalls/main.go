// Деферы: четыре вопроса, которые задают почти всегда.
//
//  1. когда вычисляются аргументы (сразу) и когда выполняется тело (при выходе);
//  2. порядок LIFO;
//  3. defer в ЦИКЛЕ — ресурсы не освобождаются до конца функции;
//  4. defer + именованный результат — единственный способ изменить возвращаемое
//     значение постфактум (этим живут recover и обертывание ошибок).
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// 1 и 2: аргументы вычисляются сразу, порядок обратный.
func orderAndArgs() {
	for i := 1; i <= 3; i++ {
		defer fmt.Println("  defer с аргументом:", i) // i вычислен сейчас: 1, 2, 3
	}
	defer func() {
		fmt.Println("  defer с замыканием видит финальное состояние")
	}()
	fmt.Println("  тело функции закончилось") // выведется ПЕРВЫМ
}

// 3. Плохо: файлы закроются только после обработки ВСЕХ путей.
func badLoop(paths []string) error {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close() // накапливается: 10 000 файлов = 10 000 открытых дескрипторов
		_ = f
	}
	return nil
}

// Хорошо: тело цикла выносится в функцию, у которой свой выход.
func goodLoop(paths []string) error {
	for _, p := range paths {
		if err := handleOne(p); err != nil {
			return err
		}
	}
	return nil
}

func handleOne(p string) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close() // закроется на выходе из handleOne, то есть на каждой итерации
	_ = f
	return nil
}

// 4. Именованный результат: defer может его прочитать и подменить.
func wrapError(path string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("обработка %s: %w", path, err)
		}
	}()
	_, err = os.Open(path)
	return err
}

// Тот же механизм превращает панику в ошибку — например, на границе библиотеки.
func safeDivide(a, b int) (res int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("восстановлено после паники: %v", r)
		}
	}()
	return a / b, nil // b == 0 → runtime error: integer divide by zero
}

// Ловушка: без ИМЕНОВАННОГО результата подменить возвращаемое значение нельзя.
func cannotChange() (res int) {
	defer func() { res *= 2 }() // сработает: res именованный
	return 21
}

func main() {
	fmt.Println("1-2. порядок и аргументы:")
	orderAndArgs()

	tmp := filepath.Join(os.TempDir(), "defer_demo.txt")
	os.WriteFile(tmp, []byte("x"), 0o600)
	defer os.Remove(tmp)

	fmt.Println("\n3. defer в цикле:")
	fmt.Println("  badLoop :", badLoop([]string{tmp, tmp}), "— дескрипторы жили до конца функции")
	fmt.Println("  goodLoop:", goodLoop([]string{tmp, tmp}), "— каждый файл закрыт на своей итерации")

	fmt.Println("\n4. defer и именованный результат:")
	err := wrapError("нет-такого-файла")
	fmt.Println("  обертка ошибки:", err)
	fmt.Println("  errors.Is(os.ErrNotExist):", errors.Is(err, os.ErrNotExist))
	fmt.Println("  safeDivide(1, 0):", func() string { _, e := safeDivide(1, 0); return e.Error() }())
	fmt.Println("  cannotChange():", cannotChange(), "← defer удвоил результат")

	// Стоимость: с Go 1.14 обычный defer инлайнится и стоит единицы наносекунд,
	// но defer в горячем цикле все равно лучше заменить явным вызовом.
}
