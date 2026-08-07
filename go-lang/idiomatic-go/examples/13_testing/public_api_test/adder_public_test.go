// Тест публичного API: единственное исключение из правила «один каталог — один пакет».
// Файл объявляет пакет pubadder_test и ИМПОРТИРУЕТ тестируемый пакет,
// поэтому видит ровно то, что видят пользователи библиотеки.
package pubadder_test

// пакет с суффиксом _test: видим только экспортируемое, как настоящий пользователь

import (
	"interviewprep/examples/13_testing/public_api_test"
	// тестируемый пакет приходится импортировать
	"testing"
)

func TestAddNumbers(t *testing.T) {
	result := pubadder.AddNumbers(2, 3)
	if result != 5 {
		t.Error("incorrect result: expected 5, got", result)
	}
}
