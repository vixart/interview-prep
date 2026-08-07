// Выравнивание структур: unsafe.Sizeof дает размер, unsafe.Offsetof — смещение поля.
// Пять структур с одинаковыми полями в разном порядке: bool,int64,bool занимает 24 байта,
// а bool,bool,int64 — всего 16. Порядок полей влияет на память.
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Sizeof(BoolInt{}),
		unsafe.Offsetof(BoolInt{}.b),
		// смещение поля от начала структуры
		unsafe.Offsetof(BoolInt{}.i))
	fmt.Println(unsafe.Sizeof(IntBool{}),
		unsafe.Offsetof(IntBool{}.i),
		unsafe.Offsetof(IntBool{}.b))
	fmt.Println()
	fmt.Println(unsafe.Sizeof(BoolIntBool{}),
		unsafe.Offsetof(BoolIntBool{}.b),
		unsafe.Offsetof(BoolIntBool{}.i),
		unsafe.Offsetof(BoolIntBool{}.b2))
	fmt.Println(unsafe.Sizeof(BoolBoolInt{}),
		unsafe.Offsetof(BoolBoolInt{}.b),
		unsafe.Offsetof(BoolBoolInt{}.b2),
		unsafe.Offsetof(BoolBoolInt{}.i))
	fmt.Println(unsafe.Sizeof(IntBoolBool{}),
		unsafe.Offsetof(IntBoolBool{}.i),
		unsafe.Offsetof(IntBoolBool{}.b),
		unsafe.Offsetof(IntBoolBool{}.b2))
}

type BoolInt struct {
	// bool + int64 = 16 байт: 1 байт данных + 7 байт выравнивания
	b bool
	i int64
}

type IntBool struct {
	i int64
	b bool
}

type BoolIntBool struct {
	// bool, int64, bool = 24 байта — два разрыва на выравнивание
	b  bool
	i  int64
	b2 bool
}

type BoolBoolInt struct {
	// те же поля, но bool'ы рядом = 16 байт. Порядок полей решает
	b  bool
	b2 bool
	i  int64
}

type IntBoolBool struct {
	i  int64
	b  bool
	b2 bool
}
