// Самый простой тест: файл *_test.go, функция TestXxx(t *testing.T),
// проверка условия и t.Error при несовпадении (t.Error продолжает тест, t.Fatal прерывает).
// Запуск: go test ./13_testing/basic_test
package adder

import "testing"

func Test_addNumbers(t *testing.T) {
	// имя обязательно начинается с Test, параметр — *testing.T
	result := addNumbers(2, 3)
	if result != 5 {
		t.Error("incorrect result: expected 5, got", result)
		// t.Error помечает тест упавшим, но продолжает его; t.Fatal прервал бы
	}
}
