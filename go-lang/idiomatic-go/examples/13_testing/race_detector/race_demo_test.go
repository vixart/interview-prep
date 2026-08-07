//go:build racedemo

// Этот тест намеренно падает (или ловится детектором гонок), поэтому он спрятан
// за тегом сборки и не ломает обычный `go test ./...`:
//
//	go test -tags racedemo -race ./13_testing/race_detector
package race

import "testing"

func TestGetCounterRace(t *testing.T) {
	if counter := getCounter(); counter != 5000 {
		// гоночная версия обычно дает МЕНЬШЕ 5000: часть инкрементов теряется
		t.Error("unexpected counter:", counter)
	}
}
