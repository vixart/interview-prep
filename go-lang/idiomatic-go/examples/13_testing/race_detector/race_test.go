// Тест безопасной версии счетчика. Гоночный вариант вынесен в race_demo_test.go
// за тег сборки, чтобы обычный go test ./... оставался зеленым.
package race

import "testing"

func TestGetCounterSafe(t *testing.T) {
	if counter := getCounterSafe(); counter != 5000 {
		// безопасная версия всегда дает ровно 5000
		t.Error("unexpected counter:", counter)
	}
}
