// Структура с неэкспортируемым полем b — обычным кодом извне оно недоступно.
package one_package

type HasUnexportedField struct {
	A int
	b bool
	// со строчной буквы = виден только внутри пакета one_package
	C string
}
