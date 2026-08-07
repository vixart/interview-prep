// Пакет с публичным API, который тестируется «снаружи» (см. adder_public_test.go).
package pubadder

func AddNumbers(x, y int) int {
	// экспортируемая функция — только такие видны из пакета pubadder_test
	return x + y
}
