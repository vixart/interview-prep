// Тестовые данные лежат в каталоге testdata: Go игнорирует его при сборке,
// а путь в тесте всегда относительный к каталогу пакета.
package text

import "testing"

func TestCountCharacters(t *testing.T) {
	total, err := CountCharacters("testdata/sample1.txt")
	// путь относительно каталога пакета; testdata Go при сборке игнорирует
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if total != 35 {
		t.Error("Expected 35, got", total)
	}
	_, err = CountCharacters("testdata/no_file.txt")
	if err == nil {
		t.Error("Expected an error")
	}
}
