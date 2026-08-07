// Разбор бинарного протокола: *(*Data)(unsafe.Pointer(&b)) кладет байты сети
// прямо на структуру, без побайтового копирования (сравни с DataFromBytes на encoding/binary).
// Обязательные оговорки: порядок байтов (isLE определяется в init) и то, что раскладка
// структуры должна ТОЧНО совпадать с форматом. Бенчмарк в main_test.go покажет разницу.
package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"unsafe"
)

type Data struct {
	Value  uint32   // 4 bytes
	Label  [10]byte // 10 bytes
	Active bool     // 1 byte
	// padded with 1 byte to make it align
}

func (d Data) String() string {
	return fmt.Sprintf(`{Value:%d, Label:%s, Active: %v}`, d.Value, string(d.Label[:]), d.Active)
}

var isLE bool

func init() {
	var x uint16 = 0xFF00
	xb := *(*[2]byte)(unsafe.Pointer(&x))
	// определяем порядок байтов машины, глядя на байты uint16
	isLE = (xb[0] == 0x00)
}

func DataFromBytesUnsafe(b [dataSize]byte) Data {
	data := *(*Data)(unsafe.Pointer(&b))
	// ВОТ ОНО: массив байт трактуется как структура, без копирования полей
	if isLE {
		data.Value = bits.ReverseBytes32(data.Value)
		// сеть — big-endian, машина обычно little-endian: разворачиваем вручную
	}
	return data
}

func DataFromBytesUnsafeSlice(b []byte) Data {
	data := *(*Data)((unsafe.Pointer)(unsafe.SliceData(b)))
	if isLE {
		data.Value = bits.ReverseBytes32(data.Value)
	}
	return data
}

func DataFromBytes(b [dataSize]byte) Data {
	d := Data{}
	d.Value = binary.BigEndian.Uint32(b[:4])
	// безопасная альтернатива: encoding/binary, поле за полем
	copy(d.Label[:], b[4:14])
	d.Active = b[14] != 0
	return d
}

const dataSize = unsafe.Sizeof(Data{})

func BytesFromDataUnsafe(d Data) [dataSize]byte {
	if isLE {
		d.Value = bits.ReverseBytes32(d.Value)
	}
	b := *(*[dataSize]byte)(unsafe.Pointer(&d))
	return b
}

func BytesFromDataUnsafeSlice(d Data) []byte {
	if isLE {
		d.Value = bits.ReverseBytes32(d.Value)
	}
	bs := unsafe.Slice((*byte)(unsafe.Pointer(&d)), dataSize)
	return bs
}

func BytesFromData(d Data) [dataSize]byte {
	out := [dataSize]byte{}
	binary.BigEndian.PutUint32(out[:4], d.Value)
	copy(out[4:14], d.Label[:])
	if d.Active {
		out[14] = 1
	}
	return out
}

func main() {
	incomingData := [dataSize]byte{0, 132, 95, 237, 80, 104, 111, 110, 101, 0, 0, 0, 0, 0, 1, 0}

	d1 := DataFromBytes(incomingData)
	d2 := DataFromBytesUnsafe(incomingData)
	if d1 != d2 {
		panic(fmt.Sprintf("%v %v", d1, d2))
	}
	fmt.Println(d1)

	b1 := BytesFromData(d1)
	b2 := BytesFromDataUnsafe(d1)
	if b1 != b2 {
		panic(fmt.Sprintf("%v %v", b1, b2))
	}
	fmt.Printf("%+v\n", b1)

	b3 := BytesFromDataUnsafeSlice(d1)
	if *(*[dataSize]byte)(b3) != b2 {
		panic(fmt.Sprintf("%v %v", b2, b3))
	}
	fmt.Printf("%+v\n", b3)

	d3 := DataFromBytesUnsafeSlice(b3)
	if d3 != d2 {
		panic(fmt.Sprintf("%v %v", d2, d3))
	}
	fmt.Println(d3)
}
