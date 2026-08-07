// cgo: C-код прямо в комментарии перед import "C" плюс соседние mylib.h/mylib.c.
// Вызов идет как C.add(3, 2). Нужен установленный компилятор C.
// Помни: cgo — про интеграцию с C-библиотеками, а не про скорость (вызовы дорогие).
package main

import "fmt"

/*
	#cgo LDFLAGS: -lm
	#include <stdio.h>
	#include <math.h>
	#include "mylib.h"

	int add(int a, int b) {
		int sum = a + b;
		printf("a: %d, b: %d, sum %d\n", a, b, sum);
		return sum;
	}
*/
import "C"

// псевдопакет C: комментарий ВЫШЕ — это и есть C-код, пустой строки между ними быть не должно

func main() {
	sum := C.add(3, 2)
	// функция из комментария выше
	fmt.Println(sum)
	fmt.Println(C.sqrt(100))
	// функция из системной библиотеки math (подключена через #cgo LDFLAGS: -lm)
	fmt.Println(C.multiply(10, 20))
	// функция из соседнего файла mylib.c — cgo собирает его сам
}
