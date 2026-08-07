// Тривиальная функция, на которой показывают t.Parallel (см. code_test.go).
package parallel

func toTest(i int) int {
	// функция без состояния — такую безопасно тестировать параллельно
	return i + 10
}
