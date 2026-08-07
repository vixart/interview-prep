// reflect + unsafe в связке: unsafe.Offsetof на неэкспортируемое поле не компилируется,
// поэтому смещение берется рефлексией (FieldByName), а запись — через unsafe.Add
// и приведение указателя. Крайняя мера, а не инструмент на каждый день.
package other_package

import (
	"fmt"
	"interviewprep/examples/14_reflect_unsafe/unsafe_unexported_field/one_package"
	"reflect"
	"unsafe"
)

func SetBUnsafe(huf *one_package.HasUnexportedField) {
	fmt.Println(unsafe.Sizeof(*huf))
	fmt.Println(unsafe.Offsetof(huf.A))
	// this line will fail to compile because you can't access unexported field b here
	//offset := unsafe.Offsetof(huf.b)
	fmt.Println(unsafe.Offsetof(huf.C))

	// use reflection to get the offset of the unexported field
	sf, _ := reflect.TypeOf(huf).Elem().FieldByName("b")
	// рефлексия видит и неэкспортируемые поля — берем смещение отсюда
	offset := sf.Offset
	fmt.Println("b offset", offset)

	// use unsafe to access the data at that position
	start := unsafe.Pointer(huf)
	pos := unsafe.Add(start, offset)
	// адресная арифметика: сдвигаемся к нужному полю
	b := (*bool)(pos)
	*b = true
	// запись в приватное поле чужого пакета — инкапсуляция в Go не защита
}
